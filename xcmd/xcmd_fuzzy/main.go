package xcmd_fuzzy

import (
	"fmt"
	"github.com/pubgo/g/xcmd"
	"github.com/pubgo/g/xerror"
	"log"

	"strings"

	"time"

	"github.com/jroimartin/gocui"
	"github.com/sahilm/fuzzy"
)

var filenamesBytes []byte
var err error

var filenames []string

var g *gocui.Gui

func Init() *xcmd.Command {

	var args = xcmd.Args(func(cmd *xcmd.Command) {

	})

	return args(&xcmd.Command{
		Use:   "fuzzy",
		Short: "fuzzy",
		RunE: func(cmd *xcmd.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			//filenamesBytes, err = ioutil.ReadFile("../testdata/ue4_filenames.txt")
			//if err != nil {
			//	panic(err)
			//}

			filenames = strings.Split(string(`UE4Client.Target.cs
UE4Editor.Target.cs
UE4Game.Target.cs
UE4Server.Target.cs
AITestSuite.Build.cs
MockAI.h
MockAI_BT.h
MockGameplayTasks.h
TestPawnAction_CallFunction.h
TestPawnAction_Log.h
TestBTDecorator_CantExecute.h
TestBTDecorator_DelayedAbort.h
TestBTTask_LatentWithFlags.h
TestBTTask_Log.h
TestBTTask_SetFlag.h
AITestSuite.cpp
AITestSuitePrivatePCH.h
TestLogger.cpp
TestPawnAction_CallFunction.cpp
TestPawnAction_Log.cpp
TestBTDecorator_CantExecute.cpp
TestBTDecorator_DelayedAbort.cpp
TestBTTask_LatentWithFlags.cpp
TestBTTask_Log.cpp
TestBTTask_SetFlag.cpp
MockAI.cpp
MockAI_BT.cpp
MockGameplayTasks.cpp
AITestsCommon.cpp
AITestsCommon.h
BBTest.cpp
BTTest.cpp
GameplayTasksTest.cpp
PawnActionsTest.cpp
ResourceIdTest.cpp
AITestSuite.h
BTBuilder.h
TestLogger.h
AllDesktopTargetPlatform.Build.cs
AllDesktopTargetPlatform.cpp
AllDesktopTargetPlatform.h
AllDesktopTargetPlatformModule.cpp
AllDesktopTargetPlatformPrivatePCH.h
AndroidDeviceDetection.Build.cs
AndroidDeviceDetectionModule.cpp
AndroidDeviceDetectionPrivatePCH.h
AndroidDeviceDetection.h
IAndroidDeviceDetection.h
IAndroidDeviceDetectionModule.h
AndroidPlatformEditor.Build.cs
AndroidPlatformEditorModule.cpp
AndroidPlatformEditorPrivatePCH.h
AndroidSDKSettings.cpp
AndroidSDKSettings.h
AndroidSDKSettingsCustomization.cpp
AndroidSDKSettingsCustomization.h
AndroidTargetSettingsCustomization.cpp
AndroidTargetSettingsCustomization.h
AndroidTargetPlatform.Build.cs
AndroidTargetDevice.h
AndroidTargetDeviceOutput.h
AndroidTargetPlatform.h
AndroidTargetPlatformModule.cpp
AndroidTargetPlatformPrivatePCH.h
Android_ASTCTargetPlatform.Build.cs
Android_ASTCTargetPlatformModule.cpp
Android_ASTCTargetPlatformPrivatePCH.h
Android_ATCTargetPlatform.Build.cs
Android_ATCTargetPlatformModule.cpp
Android_ATCTargetPlatformPrivatePCH.h
Android_DXTTargetPlatform.Build.cs
Android_DXTTargetPlatformModule.cpp
Android_DXTTargetPlatformPrivatePCH.h
Android_ETC1TargetPlatform.Build.cs
Android_ETC1TargetPlatformModule.cpp
Android_ETC1TargetPlatformPrivatePCH.h
Android_ETC2TargetPlatform.Build.cs
Android_ETC2TargetPlatformModule.cpp
Android_ETC2TargetPlatformPrivatePCH.h
Android_MultiTargetPlatform.Build.cs
Android_MultiTargetPlatformModule.cpp
Android_MultiTargetPlatformPrivatePCH.h
IAndroid_MultiTargetPlatformModule.h
Android_PVRTCTargetPlatform.Build.cs
Android_PVRTCTargetPlatformModule.cpp
Android_PVRTCTargetPlatformPrivatePCH.h
MetalShaderFormat.Build.cs
MetalBackend.cpp
MetalBackend.h
MetalShaderCompiler.cpp
MetalShaderFormat.cpp
MetalShaderFormat.h
MetalUtils.cpp
MetalUtils.h
AssetTools.Build.cs
AssetFixUpRedirectors.cpp
AssetFixUpRedirectors.h
AssetRenameManager.cpp
AssetRenameManager.h
AssetTools.cpp
AssetTools.h
AssetToolsConsoleCommands.h
AssetToolsModule.cpp
AssetToolsPrivatePCH.h
SDiscoveringAssetsDialog.cpp
SDiscoveringAssetsDialog.h
SPackageReportDialog.cpp
SPackageReportDialog.h
AssetTypeActions_AimOffset.h
AssetTypeActions_AimOffset1D.h
AssetTypeActions_AnimationAsset.cpp
AssetTypeActions_AnimationAsset.h
AssetTypeActions_AnimBlueprint.cpp
AssetTypeActions_AnimBlueprint.h
AssetTypeActions_AnimComposite.h
AssetTypeActions_AnimMontage.h
AssetTypeActions_AnimSequence.cpp
AssetTypeActions_AnimSequence.h
AssetTypeActions_BlendSpace.h
AssetTypeActions_BlendSpace1D.h
AssetTypeActions_Blueprint.cpp
AssetTypeActions_Blueprint.h
AssetTypeActions_CameraAnim.cpp
AssetTypeActions_CameraAnim.h
AssetTypeActions_Class.cpp
AssetTypeActions_Class.h
AssetTypeActions_ClassTypeBase.cpp
AssetTypeActions_ClassTypeBase.h
AssetTypeActions_CSVAssetBase.cpp
AssetTypeActions_Curve.cpp
AssetTypeActions_Curve.h
AssetTypeActions_CurveFloat.h
AssetTypeActions_CurveLinearColor.h
AssetTypeActions_CurveTable.cpp
AssetTypeActions_CurveTable.h
AssetTypeActions_CurveVector.h
AssetTypeActions_DataAsset.h
AssetTypeActions_DataTable.cpp
AssetTypeActions_DataTable.h
AssetTypeActions_DestructibleMesh.cpp
AssetTypeActions_DestructibleMesh.h
AssetTypeActions_DialogueVoice.h
AssetTypeActions_DialogueWave.cpp
AssetTypeActions_DialogueWave.h
AssetTypeActions_Enum.cpp
AssetTypeActions_Enum.h
AssetTypeActions_FbxSceneImportData.cpp
AssetTypeActions_FbxSceneImportData.h
AssetTypeActions_Font.cpp
AssetTypeActions_Font.h
AssetTypeActions_ForceFeedbackEffect.cpp
AssetTypeActions_ForceFeedbackEffect.h
AssetTypeActions_InstancedFoliageSettings.cpp
AssetTypeActions_InstancedFoliageSettings.h
AssetTypeActions_InterpData.h
AssetTypeActions_LandscapeGrassType.h
AssetTypeActions_LandscapeLayer.h
AssetTypeActions_Material.cpp
AssetTypeActions_Material.h
AssetTypeActions_MaterialFunction.cpp
AssetTypeActions_MaterialFunction.h
AssetTypeActions_MaterialInstanceConstant.cpp
AssetTypeActions_MaterialInstanceConstant.h
AssetTypeActions_MaterialInterface.cpp
AssetTypeActions_MaterialInterface.h
AssetTypeActions_MaterialParameterCollection.h
AssetTypeActions_MorphTarget.cpp
AssetTypeActions_MorphTarget.h
AssetTypeActions_NiagaraEffect.cpp
AssetTypeActions_NiagaraEffect.h
AssetTypeActions_NiagaraScript.cpp
AssetTypeActions_NiagaraScript.h
AssetTypeActions_ObjectLibrary.h
AssetTypeActions_ParticleSystem.cpp
AssetTypeActions_ParticleSystem.h
AssetTypeActions_PhysicalMaterial.cpp
AssetTypeActions_PhysicalMaterial.h
AssetTypeActions_PhysicsAsset.cpp
AssetTypeActions_PhysicsAsset.h
AssetTypeActions_ProceduralFoliageSpawner.cpp
AssetTypeActions_ProceduralFoliageSpawner.h
AssetTypeActions_Redirector.cpp
AssetTypeActions_Redirector.h
AssetTypeActions_ReverbEffect.cpp
AssetTypeActions_ReverbEffect.h
AssetTypeActions_Rig.cpp
AssetTypeActions_Rig.h
AssetTypeActions_SkeletalMesh.cpp
AssetTypeActions_SkeletalMesh.h
AssetTypeActions_Skeleton.cpp
AssetTypeActions_Skeleton.h
AssetTypeActions_SlateBrush.cpp
AssetTypeActions_SlateBrush.h
AssetTypeActions_SlateWidgetStyle.cpp
AssetTypeActions_SlateWidgetStyle.h
AssetTypeActions_SoundAttenuation.cpp
AssetTypeActions_SoundAttenuation.h
AssetTypeActions_SoundBase.cpp
AssetTypeActions_SoundBase.h
AssetTypeActions_SoundClass.cpp
AssetTypeActions_SoundClass.h
AssetTypeActions_SoundConcurrency.cpp
AssetTypeActions_SoundConcurrency.h
AssetTypeActions_SoundCue.cpp
AssetTypeActions_SoundCue.h
AssetTypeActions_SoundMix.cpp
AssetTypeActions_SoundMix.h
AssetTypeActions_SoundWave.cpp
AssetTypeActions_SoundWave.h
AssetTypeActions_StaticMesh.cpp
AssetTypeActions_StaticMesh.h
AssetTypeActions_Struct.cpp
AssetTypeActions_Struct.h
AssetTypeActions_SubsurfaceProfile.cpp
AssetTypeActions_SubsurfaceProfile.h
AssetTypeActions_Texture.cpp
AssetTypeActions_Texture.h
AssetTypeActions_Texture2D.cpp
AssetTypeActions_Texture2D.h
AssetTypeActions_TextureCube.h
AssetTypeActions_TextureLightProfile.cpp
AssetTypeActions_TextureLightProfile.h
AssetTypeActions_TextureRenderTarget.cpp
AssetTypeActions_TextureRenderTarget.h
AssetTypeActions_TextureRenderTarget2D.h
AssetTypeActions_TextureRenderTargetCube.h
AssetTypeActions_TouchInterface.cpp
AssetTypeActions_TouchInterface.h
AssetTypeActions_VectorField.h
AssetTypeActions_VectorFieldAnimated.h
AssetTypeActions_VectorFieldStatic.cpp
AssetTypeActions_VectorFieldStatic.h
AssetTypeActions_VertexAnimation.cpp
AssetTypeActions_VertexAnimation.h
AssetTypeActions_World.cpp
AssetTypeActions_World.h
AssetToolsModule.h
AssetTypeActions_Base.h
AssetTypeActions_CSVAssetBase.h
AssetTypeCategories.h
ClassTypeActions_Base.h
IAssetTools.h
IAssetTypeActions.h
IClassTypeActions.h
AudioFormatADPCM.Build.cs
AudioFormatADPCM.cpp
AudioFormatADPCM.h
AudioFormatOgg.Build.cs
AudioFormatOgg.cpp
AudioFormatOgg.h
AudioFormatOpus.Build.cs
AudioFormatOpus.cpp
AudioFormatOpus.h
AutomationController.Build.cs
AutomationCommandline.cpp
AutomationControllerManager.h
AutomationControllerManger.cpp
AutomationControllerModule.cpp
AutomationControllerModule.h
AutomationControllerPrivatePCH.h
AutomationDeviceClusterManager.cpp
AutomationDeviceClusterManager.h
AutomationReport.cpp
AutomationReport.h
AutomationReportManager.cpp
AutomationReportManager.h
AutomationController.h
IAutomationControllerManager.h
IAutomationControllerModule.h
IAutomationReport.h
AutomationWindow.Build.cs
AutomationFilter.h
AutomationPresetManager.cpp
AutomationPresetManager.h
AutomationWindowModule.cpp
AutomationWindowPrivatePCH.h
SAutomationExportMenu.cpp
SAutomationExportMenu.h
SAutomationGraphicalResultBox.cpp
SAutomationGraphicalResultBox.h
SAutomationTestItem.cpp
SAutomationTestItem.h
SAutomationTestItemContextMenu.h
SAutomationTestTreeView.h
SAutomationWindow.cpp
SAutomationWindow.h
SAutomationWindowCommandBar.cpp
SAutomationWindowCommandBar.h
AutomationWindow.h
IAutomationWindowModule.h
BlankModule.Build.cs
BlankModule.cpp
BlankModulePrivatePCH.h
BlankModule.h
BlueprintCompilerCppBackend.Build.cs
BlueprintCompilerCppBackend.cpp
BlueprintCompilerCppBackend.h
BlueprintCompilerCppBackendAnim.cpp
BlueprintCompilerCppBackendBase.cpp
BlueprintCompilerCppBackendBase.h
BlueprintCompilerCppBackendGatherDependencies.cpp
BlueprintCompilerCppBackendModule.cpp
BlueprintCompilerCppBackendModulePrivatePCH.h
BlueprintCompilerCppBackendUMG.cpp
BlueprintCompilerCppBackendUtils.cpp
BlueprintCompilerCppBackendUtils.h
BlueprintCompilerCppBackendValueHelper.cpp
BPCompilerTests.cpp
BlueprintCompilerCppBackendGatherDependencies.h
IBlueprintCompilerCppBackendModule.h
BlueprintNativeCodeGen.Build.cs
BlueprintNativeCodeGenManifest.cpp
BlueprintNativeCodeGenManifest.h
BlueprintNativeCodeGenModule.cpp
BlueprintNativeCodeGenPCH.h
BlueprintNativeCodeGenUtils.cpp
BlueprintNativeCodeGenUtils.h
NativeCodeGenCommandlineParams.cpp
NativeCodeGenCommandlineParams.h
NativeCodeGenerationTool.cpp
BlueprintNativeCodeGenModule.h
NativeCodeGenerationTool.h
BlueprintProfiler.Build.cs
BlueprintProfiler.cpp
BlueprintProfilerSupport.cpp
BlueprintProfilerPCH.h
BlueprintProfiler.h
BlueprintProfilerModule.h
BlueprintProfilerSupport.h
CollectionManager.Build.cs
Collection.cpp
Collection.h
CollectionManager.cpp
CollectionManager.h
CollectionManagerConsoleCommands.h
CollectionManagerModule.cpp
CollectionManagerPrivatePCH.h
CollectionManagerModule.h
CollectionManagerTypes.h
ICollectionManager.h
CollisionAnalyzer.Build.cs
CollisionAnalyzer.cpp
CollisionAnalyzer.h
CollisionAnalyzerModule.cpp
CollisionAnalyzerPCH.h
CollisionAnalyzerStyle.cpp
CollisionAnalyzerStyle.h
SCAQueryDetails.cpp
SCAQueryDetails.h
SCAQueryTableRow.cpp
SCAQueryTableRow.h
SCollisionAnalyzer.cpp
SCollisionAnalyzer.h
CollisionAnalyzerModule.h
ICollisionAnalyzer.h
CrashDebugHelper.Build.cs
CrashDebugHelper.cpp
CrashDebugHelperModule.cpp
CrashDebugHelperPrivatePCH.h
CrashDebugPDBCache.cpp
CrashDebugHelperLinux.cpp
CrashDebugHelperLinux.h
CrashDebugHelperMac.cpp
CrashDebugHelperMac.h
CrashDebugHelperWindows.cpp
CrashDebugHelperWindows.h
WindowsPlatformStackWalkExt.cpp
WindowsPlatformStackWalkExt.h
CrashDebugHelper.h
CrashDebugHelperModule.h
CrashDebugPDBCache.h
CrashTracker.Build.cs
AVIHandler.cpp
AVIHandler.h
CrashTrackerModule.cpp
CrashTrackerPrivatePCH.h
CrashVideoCapture.cpp
CrashVideoCapture.h
CrashTracker.h
ICrashTrackerModule.h
DerivedDataCache.Build.cs
DDCCleanup.cpp
DDCCleanup.h
DerivedDataBackendAsyncPutWrapper.h
DerivedDataBackendCorruptionWrapper.h
DerivedDataBackends.cpp
DerivedDataBackendVerifyWrapper.h
DerivedDataCache.cpp
DerivedDataLimitKeyLengthWrapper.h
FileSystemDerivedDataBackend.cpp
HierarchicalDerivedDataBackend.h
MemoryDerivedDataBackend.h
PakFileDerivedDataBackend.h
DerivedDataBackendInterface.h
DerivedDataCacheInterface.h
DerivedDataPluginInterface.h
DerivedDataUtilsInterface.h
DesktopPlatform.Build.cs
DesktopPlatformBase.cpp
DesktopPlatformBase.h
DesktopPlatformModule.cpp
DesktopPlatformPrivatePCH.h
PlatformInfo.cpp
DesktopPlatformLinux.cpp
DesktopPlatformLinux.h
DesktopPlatformMac.cpp
DesktopPlatformMac.h
MacNativeFeedbackContext.cpp
MacNativeFeedbackContext.h
DesktopPlatformWindows.cpp
DesktopPlatformWindows.h
WindowsNativeFeedbackContext.cpp
WindowsNativeFeedbackContext.h
WindowsRegistry.cpp
WindowsRegistry.h
DesktopPlatformModule.h
IDesktopPlatform.h
PlatformInfo.h
DesktopWidgets.Build.cs
DesktopWidgetsModule.cpp
DesktopWidgetsPrivatePCH.h
SFilePathPicker.cpp
SFilePathPicker.h
DeviceManager.Build.cs
DeviceManagerModule.cpp
DeviceManagerPrivatePCH.h
DeviceBrowserFilter.h
DeviceDetailsCommands.h
DeviceDetailsFeature.h
DeviceManagerModel.h
SDeviceManager.cpp
SDeviceManager.h
SDeviceApps.cpp
SDeviceApps.h
SDeviceAppsAppListRow.h
SDeviceBrowser.cpp
SDeviceBrowser.h
SDeviceBrowserContextMenu.h
SDeviceBrowserDeviceAdder.cpp
SDeviceBrowserDeviceAdder.h
SDeviceBrowserDeviceListRow.h
SDeviceBrowserFilterBar.cpp
SDeviceBrowserFilterBar.h
SDeviceBrowserTooltip.h
SDeviceDetails.cpp
SDeviceDetails.h
SDeviceDetailsFeatureListRow.h
SDeviceProcesses.cpp
SDeviceProcesses.h
SDeviceProcessesProcessListRow.h
SDeviceProcessesProcessTreeNode.h
SDeviceQuickInfo.h
SDeviceToolbar.cpp
SDeviceToolbar.h
DeviceManager.h
IDeviceManagerModule.h
DirectoryWatcher.Build.cs
DirectoryWatcherModule.cpp
DirectoryWatcherPrivatePCH.h
DirectoryWatcherTests.cpp
FileCache.cpp
FileCacheUtilities.cpp
DirectoryWatcherLinux.cpp
DirectoryWatcherLinux.h
DirectoryWatchRequestLinux.cpp
DirectoryWatchRequestLinux.h
DirectoryWatcherMac.cpp
DirectoryWatcherMac.h
DirectoryWatchRequestMac.cpp
DirectoryWatchRequestMac.h
DirectoryWatchterRunTests.cpp
DirectoryWatcherWindows.cpp
DirectoryWatcherWindows.h
DirectoryWatchRequestWindows.cpp
DirectoryWatchRequestWindows.h
DirectoryWatcherModule.h
FileCache.h
FileCacheUtilities.h
IDirectoryWatcher.h
ExternalImagePicker.Build.cs
ExternalImagePickerModule.cpp
ExternalImagePickerPrivatePCH.h
SExternalImagePicker.cpp
SExternalImagePicker.h
IExternalImagePickerModule.h
FriendsAndChat.Build.cs
FriendsAndChatModule.cpp
FriendsAndChatPrivatePCH.h
FriendsAndChatStyle.cpp
FriendsChatChromeStyle.cpp
FriendsChatStyle.cpp
FriendsComboStyle.cpp
FriendsFontStyle.cpp
FriendsListStyle.cpp
FriendsMarkupStyle.cpp
FriendsAndChat.h
FriendsAndChatStyle.h
FriendsChatChromeStyle.h
FriendsChatStyle.h
FriendsComboStyle.h
FriendsFontStyle.h
FriendsListStyle.h
FriendsMarkupStyle.h
IFriendsAndChatModule.h
FunctionalTesting.Build.cs
FuncTestRenderingComponent.h
FunctionalAITest.h
FunctionalTest.h
FunctionalTestingManager.h
FuncTestManager.cpp
FuncTestManager.h
FuncTestRenderingComponent.cpp
FunctionalAITest.cpp
FunctionalTest.cpp
FunctionalTestingManager.cpp
FunctionalTestingModule.cpp
FunctionalTestingPrivatePCH.h
ClientFuncTestPerforming.cpp
FunctionalTestingModule.h
IFuncTestManager.h
GameplayDebugger.Build.cs
GameplayDebuggerSettings.h
GameplayDebuggingComponent.h
GameplayDebuggingControllerComponent.h
GameplayDebuggingHUDComponent.h
GameplayDebuggingReplicator.h
GameplayDebuggingTypes.h
GameplayDebugger.cpp
GameplayDebuggerPrivate.h
GameplayDebuggerSettings.cpp
GameplayDebuggingComponent.cpp
GameplayDebuggingControllerComponent.cpp
GameplayDebuggingHUDComponent.cpp
GameplayDebuggingReplicator.cpp
GameplayDebugger.h
GammaUI.Build.cs
GammaUI.cpp
GammaUIPanel.cpp
GammaUIPanel.h
GammaUIPrivatePCH.h
GammaUI.h
HierarchicalLODUtilities.Build.cs
HierarchicalLODProxyProcessor.cpp
HierarchicalLODUtilities.cpp
HierarchicalLODUtilitiesModulePrivatePCH.h
HierarchicalLODProxyProcessor.h
HierarchicalLODUtilities.h
HotReload.Build.cs
HotReload.cpp
HotReloadClassReinstancer.cpp
HotReloadClassReinstancer.h
HotReloadPrivatePCH.h
IHotReload.h
HTML5PlatformEditor.Build.cs
HTML5TargetSettings.h
HTML5PlatformEditorClasses.cpp
HTML5PlatformEditorModule.cpp
HTML5PlatformEditorPrivatePCH.h
HTML5SDKSettings.cpp
HTML5SDKSettings.h
HTML5TargetPlatform.Build.cs
HTML5TargetDevice.cpp
HTML5TargetDevice.h
HTML5TargetPlatform.cpp
HTML5TargetPlatform.h
HTML5TargetPlatformModule.cpp
HTML5TargetPlatformPrivatePCH.h
IHTML5TargetPlatformModule.h
ImageCore.Build.cs
ImageCore.cpp
ImageCorePCH.h
ImageCore.h
ImageWrapper.Build.cs
BmpImageWrapper.cpp
BmpImageWrapper.h
ExrImageWrapper.cpp
ExrImageWrapper.h
IcnsImageWrapper.cpp
IcnsImageWrapper.h
IcoImageWrapper.cpp
IcoImageWrapper.h
ImageWrapperBase.cpp
ImageWrapperBase.h
ImageWrapperModule.cpp
ImageWrapperPrivatePCH.h
JpegImageWrapper.cpp
JpegImageWrapper.h
PngImageWrapper.cpp
PngImageWrapper.h
BmpImageSupport.h
ImageWrapper.h
IImageWrapper.h
IImageWrapperModule.h
InternationalizationSettings.Build.cs
InternationalizationSettingsModel.h
InternationalizationSettings.cpp
InternationalizationSettingsModel.cpp
InternationalizationSettingsModelDetails.cpp
InternationalizationSettingsModule.cpp
InternationalizationSettingsModulePrivatePCH.h
InternationalizationSettings.h
InternationalizationSettingsModelDetails.h
InternationalizationSettingsModule.h
IOSPlatformEditor.Build.cs
IOSPlatformEditorModule.cpp
IOSPlatformEditorPrivatePCH.h
IOSTargetSettingsCustomization.cpp
IOSTargetSettingsCustomization.h
SCertificateListRow.h
SProvisionListRow.h
IOSTargetPlatform.Build.cs
IOSDeviceHelper.h
IOSTargetDevice.cpp
IOSTargetDevice.h
IOSTargetPlatform.cpp
IOSTargetPlatform.h
IOSTargetPlatformModule.cpp
IOSTargetPlatformPrivatePCH.h
IOSDeviceHelperMac.cpp
IOSDeviceHelperWindows.cpp
TVOSTargetPlatform.Build.cs
TVOSTargetPlatformModule.cpp
TVOSTargetPlatformPrivatePCH.h
LauncherServices.Build.cs
LauncherServicesModule.cpp
LauncherServicesPrivatePCH.h
Launcher.cpp
Launcher.h
LauncherProjectPath.cpp
LauncherProjectPath.h
LauncherTask.h
LauncherTaskChainState.h
LauncherUATCommand.h
LauncherUATTask.h
LauncherVerifyProfileTask.h
LauncherWorker.cpp
LauncherWorker.h
LauncherDeviceGroup.h`), "\n")

			g, err = gocui.NewGui(gocui.OutputNormal)
			if err != nil {
				log.Panicln(err)
			}
			defer g.Close()

			g.Cursor = true
			g.Mouse = true

			g.SetManagerFunc(layout)

			if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
				log.Panicln(err)
			}

			if err := g.SetKeybinding("finder", gocui.KeyArrowRight, gocui.ModNone, switchToMainView); err != nil {
				log.Panicln(err)
			}

			if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, switchToSideView); err != nil {
				log.Panicln(err)
			}

			if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
				log.Panicln(err)
			}
			if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
				log.Panicln(err)
			}

			if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
				log.Panicln(err)
			}
			return
		},
	})
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func switchToSideView(g *gocui.Gui, view *gocui.View) error {
	if _, err := g.SetCurrentView("finder"); err != nil {
		return err
	}
	return nil
}

func switchToMainView(g *gocui.Gui, view *gocui.View) error {
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("finder", -1, 0, 80, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "Type pattern here. Press -> or <- to switch between panes"
		if _, err := g.SetCurrentView("finder"); err != nil {
			return err
		}
		v.Editor = gocui.EditorFunc(finder)
	}
	if v, err := g.SetView("main", 79, 0, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, "%s", filenamesBytes)
		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Title = "list of all files"
	}

	if v, err := g.SetView("results", -1, 3, 79, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Title = "Search Results"
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func finder(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		g.Update(func(gui *gocui.Gui) error {
			results, err := g.View("results")
			if err != nil {
				// handle error
			}
			results.Clear()
			t := time.Now()
			matches := fuzzy.Find(strings.TrimSpace(v.ViewBuffer()), filenames)
			elapsed := time.Since(t)
			fmt.Fprintf(results, "found %v matches in %v\n", len(matches), elapsed)
			for _, match := range matches {
				for i := 0; i < len(match.Str); i++ {
					if contains(i, match.MatchedIndexes) {
						fmt.Fprintf(results, fmt.Sprintf("\033[1m%s\033[0m", string(match.Str[i])))
					} else {
						fmt.Fprintf(results, string(match.Str[i]))
					}

				}
				fmt.Fprintln(results, "")
			}
			return nil
		})
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		g.Update(func(gui *gocui.Gui) error {
			results, err := g.View("results")
			if err != nil {
				// handle error
			}
			results.Clear()
			t := time.Now()
			matches := fuzzy.Find(strings.TrimSpace(v.ViewBuffer()), filenames)
			elapsed := time.Since(t)
			fmt.Fprintf(results, "found %v matches in %v\n", len(matches), elapsed)
			for _, match := range matches {
				for i := 0; i < len(match.Str); i++ {
					if contains(i, match.MatchedIndexes) {
						fmt.Fprintf(results, fmt.Sprintf("\033[1m%s\033[0m", string(match.Str[i])))
					} else {
						fmt.Fprintf(results, string(match.Str[i]))
					}
				}
				fmt.Fprintln(results, "")
			}
			return nil
		})
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		g.Update(func(gui *gocui.Gui) error {
			results, err := g.View("results")
			if err != nil {
				// handle error
			}
			results.Clear()
			t := time.Now()
			matches := fuzzy.Find(strings.TrimSpace(v.ViewBuffer()), filenames)
			elapsed := time.Since(t)
			fmt.Fprintf(results, "found %v matches in %v\n", len(matches), elapsed)
			for _, match := range matches {
				for i := 0; i < len(match.Str); i++ {
					if contains(i, match.MatchedIndexes) {
						fmt.Fprintf(results, fmt.Sprintf("\033[1m%s\033[0m", string(match.Str[i])))
					} else {
						fmt.Fprintf(results, string(match.Str[i]))
					}
				}
				fmt.Fprintln(results, "")
			}
			return nil
		})
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	}
}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}
