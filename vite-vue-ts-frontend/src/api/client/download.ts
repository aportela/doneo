import { axiosInstance } from "./index";

export const bgDownload = async (
  url: string,
  fileName: string = "download",
) => {
  const startTime = Date.now();

  let urlBlob: string | null = null;
  let tmpLink: HTMLAnchorElement | null = null;

  try {
    const response = await axiosInstance.get(url, {
      responseType: "blob",
    });

    const blob: Blob = response.data;

    urlBlob = URL.createObjectURL(blob);

    tmpLink = document.createElement("a");
    tmpLink.href = urlBlob;
    tmpLink.download = fileName;

    document.body.appendChild(tmpLink);
    tmpLink.click();

    return {
      fileName,
      url,
      mimeType: blob.type,
      length: blob.size,
      msTime: Date.now() - startTime,
    };
  } finally {
    if (tmpLink) {
      tmpLink.parentNode?.removeChild(tmpLink);
    }

    if (urlBlob) {
      URL.revokeObjectURL(urlBlob);
    }
  }
};
