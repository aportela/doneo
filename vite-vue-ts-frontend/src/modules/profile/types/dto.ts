export type UpdateRequest = {
  name: string;
  email: string;
  password?: string;
};

export type SaveAvatarRequest = {
  svg: string;
};

export type ProfileResponse = {
  id: string;
  name: string;
  email: string;
  createdAt: number;
  updatedAt: number | null;
};

export type EmptyResponse = {};
