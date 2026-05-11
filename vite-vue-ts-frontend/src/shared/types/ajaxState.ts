import type { APIErrorDetails } from "./apiErrorDetails";

interface AjaxStateInterface {
  ajaxRunning: boolean;
  ajaxErrors: boolean;
  ajaxErrorMessage: string | null;
  ajaxAPIErrorDetails: APIErrorDetails | null;
}

const defaultAjaxState: AjaxStateInterface = {
  ajaxRunning: false,
  ajaxErrors: false,
  ajaxErrorMessage: null,
  ajaxAPIErrorDetails: null,
};

const defaultAjaxStateRunning: AjaxStateInterface = {
  ajaxRunning: true,
  ajaxErrors: false,
  ajaxErrorMessage: null,
  ajaxAPIErrorDetails: null,
};

export { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning };
