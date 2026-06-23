/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Base URL of the Call Booking API (Prism mock in dev, real backend otherwise). */
  readonly VITE_API_BASE_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
