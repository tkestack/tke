# Expansion Framework For TKEStack

**Author**: madongdong([@MaDongdong99](https://github.com/MaDongdong99))

**Status** (20210413): Designing

## Background

[TKEStack Hook Framework](/docs/user/cluster/hooks.md) helps users to deploy the customized hooks to TKEStack, in which way user can inject customized logic during TKEStack life cycle.

There are a group of hook points for users to choose and configure. Deep users will generate their own **tke.json** with hook points specified.

Unfortunately, to make things work automatically, users usually need to generate the hook points by themselves so that it can work together with TKE-installer running.
At the same time, users also need to make logic to deal with customized materials such as Images,Charts,Static-files etc, which TKEStack is able to handle.

For automation, a **TKEStack Wrapper** is always needed to handle the scenes above, which deep users need to think about and cannot get help from TKEStack.

## Motivation

To make better use of TKEStack Hook Framework, in other words: **Do it A Native Way**, we consider **generating hook points** and **dealing with deploy materials** by a new mechanism: **Expansion Framework**.

Expansion Framework will empower users to plan hook points while they are building their own packages, let TKEStack deal with customized materials, and then get rid of things such as **TKEStack Wrapper**.

## Goals

- Users put an **Expansion Package** at a specified place, start TKE-installer, and expansion framework will do all the rest.
- Hook Points and Hook Files are auto setup.
- Customized materials are auto uploaded.
- When users disable expansion, TKEStack works as it did before. (And disable expansion means users do not put any **Expansion Package** in that place, which is also the same as what they did before)

## Non Goals

- Support multiple expansions. It's complicated to deal with conflicts between expansions, it's better to merge them together before packaging.

## Design

- Users make package layout as blow, and expansion framework does the following work:
    - Scan expansion directory to find expansion.yaml which specifies the layout
    - Verify files, hooks, images, charts materials
    - Merge hook config into tke.json
    - Copy hooks files into hook point
    - Put images, charts where they should be

![](/docs/images/expansion-framework/expansion-framework-design.png)

### How it works

- We define expansion.yaml as blow: which describes how users customize the hook framework.

```
charts: []
images: []
files: []
hooks: []
provider: []
installer_skip_steps: []
create_cluster_skip_conditions: []
create_cluster_delegate_conditions: []
```

- charts: Users put their customized charts (maybe helm charts) into the ${EXPANSION_DIR}/charts directory, and expansion framework will scan & check them, then upload to a chart repo.
- images: Users package their own images into ${EXPANSION_DIR}/images.tar.gz, and expansion framework will load/tag them, then push them to an image registry.
- files: Users put cluster hook files(for cluster.CopyFile hook) into the ${EXPANSION_DIR}/files directory, and expansion framework will scan & check them, then merge to cluster.feature.files in tke.json.
- hooks: User put install hook files(pre-install,post-cluster-ready,post-install) into the ${EXPANSION_DIR}/hooks directory, and expansion framework will scan & check them, then override origin hook files by them.
- provider: User put provider files into ${EXPANSION_DIR}/provider directory, and expansion framework will override the origin provider files
- installer_skip_steps: Install step list that users want to skip, expansion framework will merge them into config.Skipsteps in tke.json
- create_cluster_skip_conditions: Create cluster conditions list that users want to skip, expansion framework will merge them into cluster.feature.skipConditions in tke.json
- create_cluster_delegate_conditions: Create cluster conditions list that users want to delegate to an external operator, expansion framework will merge them into cluster.feature.delegateConditions(which is a new configuration for TKEStack) in tke.json
- **note**: files may need to be rendered before they are copied.

### Mechanism of **delegate_conditions**

We introduce a new mechanism called delegate_conditions when creating cluster.

This means TKEStack can delegate any condition(can be controlled by a whitelist) to an **external operator**, which will give users more abilities to customized their own TKEStack.
For delegate_conditions, TKEStack do nothing but just waiting for the condition to be **READY**, an **external operator** will take responsibility.

So an **external operator**(a beside container with tke-installer, or a sidecar with tke-platform), should complete conditions delegated to it and update cluster.Status.Conditions.

What TKEStack need to do is providing an **expansion-operator-SDK** for users to develop their own operator.

### Interactive Design

![](/docs/images/expansion-framework/expansion-framework-interactive.png)

## Detailed design

### Expansion Framework Modules

- Expansion Framework Lib - Base functions to support the framework.
- Expansion Operator SDK - If users want to develop expansion operator, give them a hand.
- Expansion Build Tools - A tool collections for users to build/check/test their expansion package easily.

![](/docs/images/expansion-framework/expansion-framework-modules.png)

### Design of Expansion Framework module

- expansion Driver module

```
type expansionDriver struct {
	log                             log.Logger
	Operator                        string `json:"operator" yaml:"operator"`
	Values                          map[string]string
	Charts                          []string `json:"charts" yaml:"charts"`
	Files                           []string `json:"files" yaml:"files"`
	Hooks                           []string `json:"hooks" yaml:"hooks"`
	Provider                        []string `json:"provider" yaml:"provider"`
	Images                          []string `json:"images" yaml:"images"`
	globalKubeconfig                []byte
	InstallerSkipSteps              []string `json:"installer_skip_steps" yaml:"installer_skip_steps"`
	CreateClusterSkipConditions     []string `json:"create_cluster_skip_conditions" yaml:"create_cluster_skip_conditions"`
	CreateClusterDelegateConditions []string `json:"create_cluster_delegate_conditions" yaml:"create_cluster_delegate_conditions"`
	// TODO: save image lists in case of installer restart to avoid load images again
	ImagesLoaded bool `json:"images_loaded" yaml:"images_loaded"`
	// if true, will prevent to pass the same expansion to cluster C
	DisablePassThroughToPlatform bool `json:"disable_pass_through_to_platform" yaml:"disable_pass_through_to_platform"`
}
```

- main method for TKEStack calling
    - scan() error
    - merge(t *TKE)
    - loadOperatorImage(ctx context.Context) error
    - patchHookFiles(ctx context.Context) error
    - startOperator(ctx context.Context) error
    - loadExpansionImages(ctx context.Context) error
    - tagExpansionImages(ctx context.Context) error
    - pushExpansionImages(ctx context.Context) error
    - uploadExpansionCharts(ctx context.Context) error

### Changes in TKEStack

- Global Changes
    - All modules with a cluster-installer (tke-installer, tke-platform-controller) add a member *expansionDriver, which supports calling of expansion functions.
    - All modules with a cluster.OnCreate method (tke-installer, tke-platform-controller) , when looping for createCluster, has break points for delegate steps, and do reload-cluster when the step fails.
- Modules
    - apiv1.ClusterFeature add member "DelegateConditions []string"
    - cluster.Interface add const "ReasonDelegated = \"Delegated\""
- TKE-installer
    - tke.do() add "mergeFeatures()"
    - initSteps add "Load Expansion Operator Image", "Patch Hook Files", "Start Expansion Operator",
    - Step::loadImages add "loadExpansionImages"
    - Step::tagImages add "tagExpansionImages"
    - Step::pushImages add "pushExpansionImages"
    - Step::writeKubeconfig add "write kubeconfig to share dir"
- tke-platform
    - TODO:
- Cluster.Create
    - Split EnsureDocker into EnsureDocker and EnsureStartDocker
    - Split EnsureKubelet into EnsureKubelet and EnsureStartKubelet


## Further Discussion

- Addon expansion: giving a **Native Way** for user to install their own apps on TKEStack.
