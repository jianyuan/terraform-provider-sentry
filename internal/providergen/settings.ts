import type { DataSource, Resource } from "./schema";
import { Glob } from "bun";

export const DATASOURCES: Array<DataSource> = Array.from(
  new Glob("./data_sources/*.ts").scanSync(),
).map((file) => import.meta.require(file).default);

export const RESOURCES: Array<Resource> = Array.from(
  new Glob("./resources/*.ts").scanSync(),
).map((file) => import.meta.require(file).default);
