"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("../_bridge");
var loginUin = 0;
var ownerUin = 0;
var loginInfo = window.LOGIN_INFO;
if (loginInfo) {
    loginUin = loginInfo.loginUin;
    ownerUin = loginInfo.ownerUin;
}
var STO_KEY = "tea_last_prj_" + loginUin + "_" + ownerUin;
function getLastProjectId() {
    var projectId = Number(localStorage.getItem(STO_KEY));
    if (!projectId) {
        // 从控制台老逻辑获取
        projectId = Number(_bridge_1._appUtil.getProjectId());
    }
    return projectId || -1;
}
exports.getLastProjectId = getLastProjectId;
function setLastProjectId(projectId) {
    localStorage.setItem(STO_KEY, String(projectId));
}
exports.setLastProjectId = setLastProjectId;
function clearLastProjectId() {
    localStorage.removeItem(STO_KEY);
}
exports.clearLastProjectId = clearLastProjectId;
function getPermitedProjectInfo() {
    return new Promise(function (resolve) {
        return _bridge_1._manager.getProjects(function (result) {
            return resolve({
                isShowAll: result.isShowAll,
                projects: (result.data || []).map(function (item) { return ({
                    projectId: item.projectId,
                    projectName: item.name,
                }); }),
            });
        });
    });
}
exports.getPermitedProjectInfo = getPermitedProjectInfo;
exports.getPermitedProjectList = function () {
    return getPermitedProjectInfo().then(function (x) { return x.projects; });
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicHJvamVjdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9icmlkZ2UvdXNlci9wcm9qZWN0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEsc0NBQWdEO0FBRWhELElBQUksUUFBUSxHQUFHLENBQUMsQ0FBQztBQUNqQixJQUFJLFFBQVEsR0FBRyxDQUFDLENBQUM7QUFFakIsSUFBTSxTQUFTLEdBQUcsTUFBTSxDQUFDLFVBQVUsQ0FBQztBQUNwQyxJQUFJLFNBQVMsRUFBRTtJQUNiLFFBQVEsR0FBRyxTQUFTLENBQUMsUUFBUSxDQUFDO0lBQzlCLFFBQVEsR0FBRyxTQUFTLENBQUMsUUFBUSxDQUFDO0NBQy9CO0FBRUQsSUFBTSxPQUFPLEdBQUcsa0JBQWdCLFFBQVEsU0FBSSxRQUFVLENBQUM7QUFFdkQsU0FBZ0IsZ0JBQWdCO0lBQzlCLElBQUksU0FBUyxHQUFXLE1BQU0sQ0FBQyxZQUFZLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxDQUFDLENBQUM7SUFFOUQsSUFBSSxDQUFDLFNBQVMsRUFBRTtRQUNkLFlBQVk7UUFDWixTQUFTLEdBQUcsTUFBTSxDQUFDLGtCQUFRLENBQUMsWUFBWSxFQUFFLENBQUMsQ0FBQztLQUM3QztJQUVELE9BQU8sU0FBUyxJQUFJLENBQUMsQ0FBQyxDQUFDO0FBQ3pCLENBQUM7QUFURCw0Q0FTQztBQUVELFNBQWdCLGdCQUFnQixDQUFDLFNBQWlCO0lBQ2hELFlBQVksQ0FBQyxPQUFPLENBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxTQUFTLENBQUMsQ0FBQyxDQUFDO0FBQ25ELENBQUM7QUFGRCw0Q0FFQztBQUVELFNBQWdCLGtCQUFrQjtJQUNoQyxZQUFZLENBQUMsVUFBVSxDQUFDLE9BQU8sQ0FBQyxDQUFDO0FBQ25DLENBQUM7QUFGRCxnREFFQztBQUVELFNBQWdCLHNCQUFzQjtJQUNwQyxPQUFPLElBQUksT0FBTyxDQUFzQixVQUFBLE9BQU87UUFDN0MsT0FBQSxrQkFBUSxDQUFDLFdBQVcsQ0FBQyxVQUFDLE1BQVc7WUFDL0IsT0FBQSxPQUFPLENBQUM7Z0JBQ04sU0FBUyxFQUFFLE1BQU0sQ0FBQyxTQUFTO2dCQUMzQixRQUFRLEVBQUUsQ0FBQyxNQUFNLENBQUMsSUFBSSxJQUFJLEVBQUUsQ0FBQyxDQUFDLEdBQUcsQ0FBQyxVQUFDLElBQVMsSUFBSyxPQUFBLENBQUM7b0JBQ2hELFNBQVMsRUFBRSxJQUFJLENBQUMsU0FBUztvQkFDekIsV0FBVyxFQUFFLElBQUksQ0FBQyxJQUFJO2lCQUN2QixDQUFDLEVBSCtDLENBRy9DLENBQUM7YUFDSixDQUFDO1FBTkYsQ0FNRSxDQUNIO0lBUkQsQ0FRQyxDQUNGLENBQUM7QUFDSixDQUFDO0FBWkQsd0RBWUM7QUFFWSxRQUFBLHNCQUFzQixHQUFHO0lBQ3BDLE9BQUEsc0JBQXNCLEVBQUUsQ0FBQyxJQUFJLENBQUMsVUFBQSxDQUFDLElBQUksT0FBQSxDQUFDLENBQUMsUUFBUSxFQUFWLENBQVUsQ0FBQztBQUE5QyxDQUE4QyxDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgX2FwcFV0aWwsIF9tYW5hZ2VyIH0gZnJvbSBcIi4uL19icmlkZ2VcIjtcblxubGV0IGxvZ2luVWluID0gMDtcbmxldCBvd25lclVpbiA9IDA7XG5cbmNvbnN0IGxvZ2luSW5mbyA9IHdpbmRvdy5MT0dJTl9JTkZPO1xuaWYgKGxvZ2luSW5mbykge1xuICBsb2dpblVpbiA9IGxvZ2luSW5mby5sb2dpblVpbjtcbiAgb3duZXJVaW4gPSBsb2dpbkluZm8ub3duZXJVaW47XG59XG5cbmNvbnN0IFNUT19LRVkgPSBgdGVhX2xhc3RfcHJqXyR7bG9naW5VaW59XyR7b3duZXJVaW59YDtcblxuZXhwb3J0IGZ1bmN0aW9uIGdldExhc3RQcm9qZWN0SWQoKSB7XG4gIGxldCBwcm9qZWN0SWQ6IG51bWJlciA9IE51bWJlcihsb2NhbFN0b3JhZ2UuZ2V0SXRlbShTVE9fS0VZKSk7XG5cbiAgaWYgKCFwcm9qZWN0SWQpIHtcbiAgICAvLyDku47mjqfliLblj7DogIHpgLvovpHojrflj5ZcbiAgICBwcm9qZWN0SWQgPSBOdW1iZXIoX2FwcFV0aWwuZ2V0UHJvamVjdElkKCkpO1xuICB9XG5cbiAgcmV0dXJuIHByb2plY3RJZCB8fCAtMTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHNldExhc3RQcm9qZWN0SWQocHJvamVjdElkOiBudW1iZXIpIHtcbiAgbG9jYWxTdG9yYWdlLnNldEl0ZW0oU1RPX0tFWSwgU3RyaW5nKHByb2plY3RJZCkpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gY2xlYXJMYXN0UHJvamVjdElkKCkge1xuICBsb2NhbFN0b3JhZ2UucmVtb3ZlSXRlbShTVE9fS0VZKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGdldFBlcm1pdGVkUHJvamVjdEluZm8oKSB7XG4gIHJldHVybiBuZXcgUHJvbWlzZTxQZXJtaXRlZFByb2plY3RJbmZvPihyZXNvbHZlID0+XG4gICAgX21hbmFnZXIuZ2V0UHJvamVjdHMoKHJlc3VsdDogYW55KSA9PlxuICAgICAgcmVzb2x2ZSh7XG4gICAgICAgIGlzU2hvd0FsbDogcmVzdWx0LmlzU2hvd0FsbCxcbiAgICAgICAgcHJvamVjdHM6IChyZXN1bHQuZGF0YSB8fCBbXSkubWFwKChpdGVtOiBhbnkpID0+ICh7XG4gICAgICAgICAgcHJvamVjdElkOiBpdGVtLnByb2plY3RJZCxcbiAgICAgICAgICBwcm9qZWN0TmFtZTogaXRlbS5uYW1lLFxuICAgICAgICB9KSksXG4gICAgICB9KVxuICAgIClcbiAgKTtcbn1cblxuZXhwb3J0IGNvbnN0IGdldFBlcm1pdGVkUHJvamVjdExpc3QgPSAoKSA9PlxuICBnZXRQZXJtaXRlZFByb2plY3RJbmZvKCkudGhlbih4ID0+IHgucHJvamVjdHMpO1xuXG5leHBvcnQgaW50ZXJmYWNlIFBlcm1pdGVkUHJvamVjdEluZm8ge1xuICAvKipcbiAgICog5b2T5YmN55So5oi35piv5ZCm5pyJ5YW35aSH5p+l55yL5omA5pyJ6aG555uu55qE5p2D6ZmQXG4gICAqL1xuICBpc1Nob3dBbGw6IGJvb2xlYW47XG5cbiAgLyoqXG4gICAqIOW9k+WJjeeUqOaIt+acieadg+mZkOeahOmhueebruWIl+ihqFxuICAgKi9cbiAgcHJvamVjdHM6IFByb2plY3RJdGVtW107XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgUHJvamVjdEl0ZW0ge1xuICAvKipcbiAgICog6aG555uuIElE77yM5Li6IDAg6KGo56S66buY6K6k6aG555uuXG4gICAqL1xuICBwcm9qZWN0SWQ6IG51bWJlcjtcblxuICAvKipcbiAgICog6aG555uu5ZCN56ewXG4gICAqL1xuICBwcm9qZWN0TmFtZTogc3RyaW5nO1xufVxuIl19