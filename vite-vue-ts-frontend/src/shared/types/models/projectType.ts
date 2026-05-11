interface ProjectTypeInterface {
  id: string;
  name: string;
  hexColor: string;
}

class ProjectTypeClass implements ProjectTypeInterface {
  id: string;
  name: string;
  hexColor: string;
  constructor(item: ProjectTypeInterface) {
    this.id = item.id;
    this.name = item.name;
    this.hexColor = item.hexColor;
  }
}

const maxNameLength = 32;

export { type ProjectTypeInterface, ProjectTypeClass, maxNameLength };
