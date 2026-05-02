interface ProjectTypeInterface {
  id: string;
  name: string;
  index: number;
  hexColor: string;
}

class ProjectTypeClass implements ProjectTypeInterface {
  id: string;
  name: string;
  index: number;
  hexColor: string;
  constructor(item: ProjectTypeInterface) {
    this.id = item.id;
    this.name = item.name;
    this.index = item.index;
    this.hexColor = item.hexColor;
  }
}

export { type ProjectTypeInterface, type ProjectTypeClass };
