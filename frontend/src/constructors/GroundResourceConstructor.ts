import type { GroundResource } from "../types/GroundResource";

export const createDefaultGroundResource = (overrides: Partial<GroundResource> = {}): GroundResource => ({
  id: 0,
  resource_code: "",
  resource_type: "GPU",
  location: "",
  availability_status: "AVAILABLE",
  maintenance_due_at: "",
  owner_team: "",
  ...overrides
});

export const createGroundResourceForm = createDefaultGroundResource;
export const createGroundResourceResponse = createDefaultGroundResource;
