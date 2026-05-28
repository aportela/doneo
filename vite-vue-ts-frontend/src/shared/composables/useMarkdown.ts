import MarkdownIt from "markdown-it";
import DOMPurify from "dompurify";

// @ts-expect-error
import TurndownService from "turndown";

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  breaks: true,
});

const turndown = new TurndownService();

export function useMarkdown() {
  const render = (text: string) => DOMPurify.sanitize(md.render(text));

  const toMarkdown = (html: string) => turndown.turndown(html);

  return {
    md,
    turndown,
    render,
    toMarkdown,
  };
}
