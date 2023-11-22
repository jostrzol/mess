"use client";

import { RoomApi } from "@/api/room";
import { RulesApi } from "@/api/rules";
import { Button } from "@/components/form/button";
import { Input } from "@/components/form/input";
import { UploadFile } from "@/components/form/uploadFile";
import { Navbar } from "@/components/navbar";
import { useMessApi } from "@/contexts/messApiContext";
import { useTheme } from "@/contexts/themeContext";
import { Editor } from "@monaco-editor/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { redirect } from "next/navigation";
import { useEffect, useState } from "react";
import { RoomPageParams } from "../layout";

const RoomPage = ({ params: { roomId } }: RoomPageParams) => {
  const [filename, setFilename] = useState<string>("");
  const [isDirty, setIsDirty] = useState(false);

  // intentionally useState instead of useRef to trigger
  // rerender on change
  const [editor, setEditor] = useState<any | null>(null);
  const { theme } = useTheme();

  const roomApi = useMessApi(RoomApi);
  const rulesApi = useMessApi(RulesApi);

  const client = useQueryClient();
  const { data: room, isSuccess } = useQuery({
    queryKey: ["room", roomId],
    queryFn: () => roomApi.getRoom(roomId),
  });
  const { data: file } = useQuery({
    queryKey: ["room", roomId, "rules"],
    queryFn: () => roomApi.getRules(roomId),
  });
  const { mutateAsync: format } = useMutation({
    mutationFn: (src: string) => rulesApi.format(src),
  });
  const { mutate: save } = useMutation({
    mutationFn: () => roomApi.saveRules(roomId, filename, editor.getValue()),
    onSuccess: () => {
      setIsDirty(false);
      client.setQueryData(["room", roomId], {
        ...room,
        rulesFilename: filename,
      });
    },
  });

  useEffect(() => {
    if (room !== undefined) {
      setFilename(room.rulesFilename);
      setIsDirty(false);
    }
  }, [room]);

  if (!isSuccess) {
    return null;
  }
  if (room.isStarted) {
    redirect(`/rooms/${room.id}/game`);
  }

  return (
    <>
      <Navbar>
        <form
          className="w-full flex gap-2"
          onSubmit={(e) => {
            e.preventDefault();
            save();
          }}
        >
          <Input
            required
            className="font-mono placeholder:font-sans grow"
            placeholder="Enter filename"
            spellCheck={false}
            value={filename}
            onChange={(e) => setFilename(e.currentTarget.value)}
          />
          <UploadFile
            onChange={async (e) => {
              const file = e.currentTarget.files?.[0];
              if (file) {
                setFilename(file.name);
                editor.setValue(await file.text());
                save();
              }
            }}
          >
            Upload
          </UploadFile>
          <Button
            type="button"
            onClick={() =>
              editor?.getAction("editor.action.formatDocument").run()
            }
          >
            Format
          </Button>
          <Button
            type="submit"
            disabled={editor === null || filename === "" || !isDirty}
          >
            Save
          </Button>
        </form>
      </Navbar>
      <Editor
        defaultLanguage="hcl"
        theme={theme.editor}
        value={file}
        defaultPath={room.rulesFilename}
        onChange={() => setIsDirty(true)}
        onMount={(mountedEditor, monaco) => {
          setEditor(mountedEditor);
          monaco.languages.registerDocumentFormattingEditProvider("hcl", {
            provideDocumentFormattingEdits: async (model: any) => {
              const value = model.getValue(undefined, undefined);
              const out = await format(value);
              const range = model.getFullModelRange();
              return [{ range, text: out }];
            },
          });
        }}
      />
    </>
  );
};

export default RoomPage;
