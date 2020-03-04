declare namespace nmc {
  type RegionId = string | number;
  type ProjectId = string | number;
  interface AppUtil {
    getRegionId: () => RegionId;
    setRegionId: (regionId: RegionId) => void;
    getProjectId: (skipAll?: boolean) => ProjectId;
    setProjectId: (ProjectId: ProjectId) => void;
  }
}
