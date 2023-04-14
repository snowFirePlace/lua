local LrApplication = import 'LrApplication'
local LrDialogs = import 'LrDialogs'
local LrTasks = import 'LrTasks'
local LrPathUtils = import 'LrPathUtils'

LrTasks.startAsyncTask(function ()
  local catalog = LrApplication.activeCatalog()

  local photo = catalog:getTargetPhoto()
  if photo == nil then
    LrDialogs.message("Image search YImages", "Please select a photo")
    return
  end
  local path = photo:getRawMetadata("path")

  local command
	local quotedCommand
  
  command = '"' .. LrPathUtils.child( LrPathUtils.child( _PLUGIN.path, "imgur" ), "imgur.exe" ) .. '" ' .. '"' .. path .. '" ' .. '"'
  quotedCommand = '"' .. command .. '"'
  local msg = string.format("%q", '"'.."asf")
  if LrTasks.execute( quotedCommand ) ~= 0 then
    LrDialogs.message("First way", msg)
  end
  LrDialogs.message("Hello World", msg)
end)

