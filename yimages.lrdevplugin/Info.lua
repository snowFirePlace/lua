return {
  VERSION = { major=1, minor=0, revision=0, },

  LrSdkVersion = 9.0,
  LrSdkMinimumVersion = 4.0,

	LrToolkitIdentifier = 'com.adobe.lightroom.sdk.export.creator',
  LrPluginName = "Image search using YImages",
  LrPluginInfoUrl="https://github.com/snowFirePlace/",
  LrLibraryMenuItems = {
    {
      title = "Search image",
      file = "search.lua",
    },
  },  
}