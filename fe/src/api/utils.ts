type ErrorDto = {
  Status: number;
  Message: string;
  Validation?: {
    Field: string;
    Message: string;
  }[];
}

export const throwIfError = async (res: Response): Promise<void> => {
  if (!res.ok) {
    let msg = "failed to fetch data"
    try {
      const json: ErrorDto = await res.json()
      msg += ": " + json.Message
    } finally {
      throw new Error(msg);
    }
  }
}
