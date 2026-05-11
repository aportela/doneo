interface ProjectStatusInterface {
  id: string;
  name: string;
  hexColor: string;
  index: number;
}

class ProjectStatusClass implements ProjectStatusInterface {
  id: string;
  name: string;
  hexColor: string;
  index: number;
  constructor(item: ProjectStatusInterface) {
    this.id = item.id;
    this.name = item.name;
    this.hexColor = item.hexColor;
    this.index = item.index;
  }
}

const maxNameLength = 16;

export { type ProjectStatusInterface, ProjectStatusClass, maxNameLength };
