import { defineComponent, type PropType, type VNodeChild } from "vue";

export default defineComponent({
  name: "RenderCell",

  props: {
    render: {
      type: Function as PropType<(row: any) => VNodeChild>,
      required: true,
    },
    row: {
      type: Object as PropType<any>,
      required: true,
    },
  },

  setup(props) {
    return () => props.render(props.row);
  },
});
