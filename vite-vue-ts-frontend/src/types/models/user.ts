interface UserInterface {
  id: string;
  name: string;
  password?: string;
  email: string;
  isSuperUser: boolean;
  createdAt: number;
  updatedAt: number;
  deletedAt: number;
  avatar: string;
}

class UserClass implements UserInterface {
  id: string;
  name: string;
  password?: string;
  email: string;
  isSuperUser: boolean;
  createdAt: number;
  updatedAt: number;
  deletedAt: number;
  avatar: string;

  constructor(item: UserInterface) {
    this.id = item.id;
    this.name = item.name;
    this.password = item.password;
    this.email = item.email;
    this.isSuperUser = item.isSuperUser;
    this.createdAt = item.createdAt;
    this.updatedAt = item.updatedAt;
    this.deletedAt = item.deletedAt;
    this.avatar = item.avatar;
  }
}

const maxNameLength = 32;
const maxEmailLength = 255;

export { type UserInterface, UserClass, maxNameLength, maxEmailLength };
