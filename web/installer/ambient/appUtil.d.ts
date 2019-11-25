declare namespace nmc {
    interface AppUtil{
        getRegionId: ()=>RegionId;
        setRegionId: (regionId: RegionId)=>void;
        getProjectId: (skipAll?: boolean)=>ProjectId;
        setProjectId: (ProjectId: ProjectId)=>void;
        submitForm: (url: string, data: any, method?: string)=>void;
    }
}
