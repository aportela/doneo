const hasExtension = (filename: string, exts: string[]): boolean => {
  return new RegExp(`\\.(${exts.join("|")})$`, "i").test(filename);
};

const allowPreview = (f: string): boolean =>
  hasExtension(f, [
    "jpg",
    "jpeg",
    "png",
    "gif",
    "svg",
    "webp",
    "mp3",
    "ogg",
    "pdf",
  ]);

const isImage = (f: string): boolean =>
  hasExtension(f, ["jpg", "jpeg", "png", "gif", "svg", "webp"]);

const isAudio = (f: string): boolean => hasExtension(f, ["mp3", "ogg"]);

const isPDF = (f: string): boolean => hasExtension(f, ["pdf"]);

export { allowPreview, isImage, isAudio, isPDF };
