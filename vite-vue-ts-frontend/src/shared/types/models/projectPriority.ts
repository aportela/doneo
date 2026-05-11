interface ProjectPriorityInterface {
  id: string;
  name: string;
  hexColor: string;
  index: number;
}

class ProjectPriorityClass implements ProjectPriorityInterface {
  id: string;
  name: string;
  hexColor: string;
  index: number;
  constructor(item: ProjectPriorityInterface) {
    this.id = item.id;
    this.name = item.name;
    this.hexColor = item.hexColor;
    this.index = item.index;
  }
}

const maxNameLength = 16;

export { type ProjectPriorityInterface, ProjectPriorityClass, maxNameLength };
