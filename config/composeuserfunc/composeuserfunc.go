// This package is strongly based on github.com/hashicorp/hcl/v2/ext/userfunc.
// Unfortunately its source code had to be copied here, as it isn't extensible at all
package composeuserfunc

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/userfunc"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var compositeFuncBodySchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "params",
			Required: true,
		},
		{
			Name:     "variadic_param",
			Required: false,
		},
		{
			Name:     "result",
			Required: true,
		},
	},
}

func DecodeCompositeUserFunctions(body hcl.Body, blockType string, contextFunc userfunc.ContextFunc) (funcs map[string]function.Function, remain hcl.Body, diags hcl.Diagnostics) {
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       blockType,
				LabelNames: []string{"name"},
			},
		},
	}

	content, remain, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return nil, remain, diags
	}

	// first call to getBaseCtx will populate context, and then the same
	// context will be used for all subsequent calls. It's assumed that
	// all functions in a given body should see an identical context.
	var baseCtx *hcl.EvalContext
	getBaseCtx := func() *hcl.EvalContext {
		if baseCtx == nil {
			if contextFunc != nil {
				baseCtx = contextFunc()
			}
		}
		// baseCtx might still be nil here, and that's okay
		return baseCtx
	}

	funcs = make(map[string]function.Function)
Blocks:
	for _, block := range content.Blocks {
		name := block.Labels[0]
		funcContent, funcDiags := block.Body.Content(compositeFuncBodySchema)
		diags = append(diags, funcDiags...)
		if funcDiags.HasErrors() {
			continue
		}

		paramsExpr := funcContent.Attributes["params"].Expr
		resultExpr := funcContent.Attributes["result"].Expr
		var varParamExpr hcl.Expression
		if funcContent.Attributes["variadic_param"] != nil {
			varParamExpr = funcContent.Attributes["variadic_param"].Expr
		}

		var params []string
		var varParam string

		paramExprs, paramsDiags := hcl.ExprList(paramsExpr)
		diags = append(diags, paramsDiags...)
		if paramsDiags.HasErrors() {
			continue
		}
		for _, paramExpr := range paramExprs {
			param := hcl.ExprAsKeyword(paramExpr)
			if param == "" {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid param element",
					Detail:   "Each parameter name must be an identifier.",
					Subject:  paramExpr.Range().Ptr(),
				})
				continue Blocks
			}
			params = append(params, param)
		}

		if varParamExpr != nil {
			varParam = hcl.ExprAsKeyword(varParamExpr)
			if varParam == "" {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid variadic_param",
					Detail:   "The variadic parameter name must be an identifier.",
					Subject:  varParamExpr.Range().Ptr(),
				})
				continue
			}
		}

		spec := &function.Spec{}
		for _, paramName := range params {
			spec.Params = append(spec.Params, function.Parameter{
				Name: paramName,
				Type: cty.DynamicPseudoType,
			})
		}
		if varParamExpr != nil {
			spec.VarParam = &function.Parameter{
				Name: varParam,
				Type: cty.DynamicPseudoType,
			}
		}
		impl := func(args []cty.Value) (cty.Value, error) {
			ctx := getBaseCtx()
			ctx = ctx.NewChild()
			ctx.Variables = make(map[string]cty.Value)

			// The cty function machinery guarantees that we have at least
			// enough args to fill all of our params.
			for i, paramName := range params {
				ctx.Variables[paramName] = args[i]
			}
			if spec.VarParam != nil {
				varArgs := args[len(params):]
				ctx.Variables[varParam] = cty.TupleVal(varArgs)
			}

			var diags hcl.Diagnostics
			steps, diags := hcl.ExprMap(resultExpr)
			if diags.HasErrors() {
				return cty.DynamicVal, diags
			}

			var returnExpr hcl.Expression
			for _, step := range steps {
				varNameValue, keyDiags := step.Key.Value(nil)
				diags = diags.Extend(keyDiags)
				if diags.HasErrors() {
					continue
				}

				varName := varNameValue.AsString()
				if varName == "return" {
					returnExpr = step.Value
					continue
				}

				varValue, valueDiags := step.Value.Value(ctx)
				diags = diags.Extend(valueDiags)
				if diags.HasErrors() {
					continue
				}

				ctx.Variables[varName] = varValue
			}

			if diags.HasErrors() {
				return cty.DynamicVal, diags
			}

			result, diags := returnExpr.Value(ctx)
			if diags.HasErrors() {
				return cty.DynamicVal, diags
			}

			return result, nil
		}
		spec.Type = function.StaticReturnType(cty.DynamicPseudoType)
		spec.Impl = func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			return impl(args)
		}
		funcs[name] = function.New(spec)
	}

	return funcs, remain, diags
}
