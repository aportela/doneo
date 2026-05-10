export type User = {
  id: string;
  name: string;
  email: string;
  isSuperUser: boolean;
  createdAt: number;
  updatedAt: number | null;
  deletedAt: number | null;
  avatarURL: string;
};

export type AddResponseInterface = {
  user: User;
};

export type UpdateResponseInterface = {
  user: User;
};

export type GetResponseInterface = {
  user: User;
};

export type SearchResponseInterface = {
  users: User[];
};
