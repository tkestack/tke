"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("../_bridge");
var event_1 = require("./event");
var $ = _bridge_1._bridge("$");
var _login = _bridge_1._bridge("login");
/**
 * 清除用户当前登录态，并弹出登录对话框
 */
function login() {
    return _login.show();
}
exports.login = login;
/**
 * 清除用户当前登录态
 */
function logout() {
    return _login.logout();
}
exports.logout = logout;
if ($) {
    /**
     * 用户登录态失效时触发
     */
    $(document).on("logout", function () {
        return event_1.userEmitter.emit("invalidate", { source: "logout" });
    });
    $(document).on("accountChanged", function () {
        return event_1.userEmitter.emit("invalidate", { source: "accountChanged" });
    });
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibG9naW4uanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvYnJpZGdlL3VzZXIvbG9naW4udHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSxzQ0FBcUM7QUFDckMsaUNBQXNDO0FBRXRDLElBQU0sQ0FBQyxHQUFHLGlCQUFPLENBQUMsR0FBRyxDQUFDLENBQUM7QUFDdkIsSUFBTSxNQUFNLEdBQUcsaUJBQU8sQ0FBQyxPQUFPLENBQUMsQ0FBQztBQUVoQzs7R0FFRztBQUNILFNBQWdCLEtBQUs7SUFDbkIsT0FBTyxNQUFNLENBQUMsSUFBSSxFQUFFLENBQUM7QUFDdkIsQ0FBQztBQUZELHNCQUVDO0FBRUQ7O0dBRUc7QUFDSCxTQUFnQixNQUFNO0lBQ3BCLE9BQU8sTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO0FBQ3pCLENBQUM7QUFGRCx3QkFFQztBQUVELElBQUksQ0FBQyxFQUFFO0lBQ0w7O09BRUc7SUFDSCxDQUFDLENBQUMsUUFBUSxDQUFDLENBQUMsRUFBRSxDQUFDLFFBQVEsRUFBRTtRQUN2QixPQUFBLG1CQUFXLENBQUMsSUFBSSxDQUFDLFlBQVksRUFBRSxFQUFFLE1BQU0sRUFBRSxRQUFRLEVBQUUsQ0FBQztJQUFwRCxDQUFvRCxDQUNyRCxDQUFDO0lBQ0YsQ0FBQyxDQUFDLFFBQVEsQ0FBQyxDQUFDLEVBQUUsQ0FBQyxnQkFBZ0IsRUFBRTtRQUMvQixPQUFBLG1CQUFXLENBQUMsSUFBSSxDQUFDLFlBQVksRUFBRSxFQUFFLE1BQU0sRUFBRSxnQkFBZ0IsRUFBRSxDQUFDO0lBQTVELENBQTRELENBQzdELENBQUM7Q0FDSCIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IF9icmlkZ2UgfSBmcm9tIFwiLi4vX2JyaWRnZVwiO1xuaW1wb3J0IHsgdXNlckVtaXR0ZXIgfSBmcm9tIFwiLi9ldmVudFwiO1xuXG5jb25zdCAkID0gX2JyaWRnZShcIiRcIik7XG5jb25zdCBfbG9naW4gPSBfYnJpZGdlKFwibG9naW5cIik7XG5cbi8qKlxuICog5riF6Zmk55So5oi35b2T5YmN55m75b2V5oCB77yM5bm25by55Ye655m75b2V5a+56K+d5qGGXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBsb2dpbigpIHtcbiAgcmV0dXJuIF9sb2dpbi5zaG93KCk7XG59XG5cbi8qKlxuICog5riF6Zmk55So5oi35b2T5YmN55m75b2V5oCBXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBsb2dvdXQoKSB7XG4gIHJldHVybiBfbG9naW4ubG9nb3V0KCk7XG59XG5cbmlmICgkKSB7XG4gIC8qKlxuICAgKiDnlKjmiLfnmbvlvZXmgIHlpLHmlYjml7bop6blj5FcbiAgICovXG4gICQoZG9jdW1lbnQpLm9uKFwibG9nb3V0XCIsICgpID0+XG4gICAgdXNlckVtaXR0ZXIuZW1pdChcImludmFsaWRhdGVcIiwgeyBzb3VyY2U6IFwibG9nb3V0XCIgfSlcbiAgKTtcbiAgJChkb2N1bWVudCkub24oXCJhY2NvdW50Q2hhbmdlZFwiLCAoKSA9PlxuICAgIHVzZXJFbWl0dGVyLmVtaXQoXCJpbnZhbGlkYXRlXCIsIHsgc291cmNlOiBcImFjY291bnRDaGFuZ2VkXCIgfSlcbiAgKTtcbn1cbiJdfQ==