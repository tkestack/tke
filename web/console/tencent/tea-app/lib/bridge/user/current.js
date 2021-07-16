"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("../_bridge");
// 当前用户数据，缓存第一次拉取
var currentUserData;
exports.current = function () {
    return currentUserData ||
        (currentUserData = new Promise(function (resolve) {
            return _bridge_1._manager.getComData(function (data) {
                if (!data) {
                    resolve(null);
                    return;
                }
                var userInfo = data.userInfo, ownerInfo = data.ownerInfo, appId = data.appId;
                var loginUin = Number(userInfo.uin);
                var ownerUin = Number(userInfo.ownerUin);
                var identity = null;
                if (ownerInfo && ownerInfo.authDetail && "type" in ownerInfo.authDetail) {
                    var _a = ownerInfo.authDetail, type = _a.type, authType = _a.authType, area = _a.area;
                    // 字段解释：http://tapd.oa.com/QCloud_2015/markdown_wikis/#1010103951008322599
                    identity = {
                        subjectType: +type,
                        authType: +authType,
                        authArea: +area,
                    };
                }
                resolve({
                    isOwner: loginUin === ownerUin,
                    loginUin: loginUin,
                    ownerUin: ownerUin,
                    appId: +appId,
                    identity: identity,
                    nickName: userInfo.nick,
                    displayName: userInfo.displayName,
                });
            });
        }));
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY3VycmVudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9icmlkZ2UvdXNlci9jdXJyZW50LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEsc0NBQXNDO0FBK0R0QyxpQkFBaUI7QUFDakIsSUFBSSxlQUFxQyxDQUFDO0FBQzdCLFFBQUEsT0FBTyxHQUFHO0lBQ3JCLE9BQUEsZUFBZTtRQUNmLENBQUMsZUFBZSxHQUFHLElBQUksT0FBTyxDQUFDLFVBQUEsT0FBTztZQUNwQyxPQUFBLGtCQUFRLENBQUMsVUFBVSxDQUFDLFVBQUMsSUFBUztnQkFDNUIsSUFBSSxDQUFDLElBQUksRUFBRTtvQkFDVCxPQUFPLENBQUMsSUFBSSxDQUFDLENBQUM7b0JBQ2QsT0FBTztpQkFDUjtnQkFDTyxJQUFBLHdCQUFRLEVBQUUsMEJBQVMsRUFBRSxrQkFBSyxDQUFVO2dCQUU1QyxJQUFNLFFBQVEsR0FBRyxNQUFNLENBQUMsUUFBUSxDQUFDLEdBQUcsQ0FBQyxDQUFDO2dCQUN0QyxJQUFNLFFBQVEsR0FBRyxNQUFNLENBQUMsUUFBUSxDQUFDLFFBQVEsQ0FBQyxDQUFDO2dCQUUzQyxJQUFJLFFBQVEsR0FBK0IsSUFBSSxDQUFDO2dCQUVoRCxJQUFJLFNBQVMsSUFBSSxTQUFTLENBQUMsVUFBVSxJQUFJLE1BQU0sSUFBSSxTQUFTLENBQUMsVUFBVSxFQUFFO29CQUNqRSxJQUFBLHlCQUErQyxFQUE3QyxjQUFJLEVBQUUsc0JBQVEsRUFBRSxjQUE2QixDQUFDO29CQUN0RCwwRUFBMEU7b0JBQzFFLFFBQVEsR0FBRzt3QkFDVCxXQUFXLEVBQUUsQ0FBQyxJQUFJO3dCQUNsQixRQUFRLEVBQUUsQ0FBQyxRQUFRO3dCQUNuQixRQUFRLEVBQUUsQ0FBQyxJQUFJO3FCQUNoQixDQUFDO2lCQUNIO2dCQUVELE9BQU8sQ0FBQztvQkFDTixPQUFPLEVBQUUsUUFBUSxLQUFLLFFBQVE7b0JBQzlCLFFBQVEsVUFBQTtvQkFDUixRQUFRLFVBQUE7b0JBQ1IsS0FBSyxFQUFFLENBQUMsS0FBSztvQkFDYixRQUFRLFVBQUE7b0JBQ1IsUUFBUSxFQUFFLFFBQVEsQ0FBQyxJQUFJO29CQUN2QixXQUFXLEVBQUUsUUFBUSxDQUFDLFdBQVc7aUJBQ2xDLENBQUMsQ0FBQztZQUNMLENBQUMsQ0FBQztRQS9CRixDQStCRSxDQUNILENBQUM7QUFsQ0YsQ0FrQ0UsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IF9tYW5hZ2VyIH0gZnJvbSBcIi4uL19icmlkZ2VcIjtcblxuZXhwb3J0IGludGVyZmFjZSBBcHBVc2VyRGF0YSB7XG4gIC8qKiDmmK/lkKbkuLrkuLvotKblj7cgKi9cbiAgaXNPd25lcjogYm9vbGVhbjtcblxuICAvKiog5b2T5YmN55So5oi355m75b2V55qEIFVJTiAqL1xuICBsb2dpblVpbjogbnVtYmVyO1xuXG4gIC8qKiDlvZPliY3nlKjmiLfnmbvlvZXnmoTkuLvotKblj7cgVUlOICovXG4gIG93bmVyVWluOiBudW1iZXI7XG5cbiAgLyoqIOW9k+WJjeeUqOaIt+eZu+W9leeahOS4u+i0puWPtyBBUFBJRCAqL1xuICBhcHBJZDogbnVtYmVyO1xuXG4gIC8qKiDlvZPliY3nlKjmiLfnmoTlrp7lkI3orqTor4Hkv6Hmga/vvIzlpoLmnpzmnKrlrp7lkI3orqTor4HvvIzmraTlrZfmrrXkuLogbnVsbCAqL1xuICBpZGVudGl0eTogQXBwVXNlcklkZW50aXR5SW5mbyB8IG51bGw7XG5cbiAgLyoqIOeUqOaIt+aYteensCAqL1xuICBuaWNrTmFtZT86IHN0cmluZztcblxuICAvKiog55So5oi35qCH6K+G5ZCN56ew77yM5ZCr5byA5Y+R5ZWG5L+h5oGvICovXG4gIGRpc3BsYXlOYW1lPzogc3RyaW5nO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEFwcFVzZXJJZGVudGl0eUluZm8ge1xuICAvKipcbiAgICog6K6k6K+B5Li75L2T57G75Z6LXG4gICAqICAgLSBgMGA6IOS4quS6ulxuICAgKiAgIC0gYDFgOiDkvIHkuJpcbiAgICovXG4gIHN1YmplY3RUeXBlOiBudW1iZXI7XG5cbiAgLyoqXG4gICAqIOiupOivgea4oOmBk1xuICAgKiAgIC0gYDBgOiDmnKrnn6VcbiAgICogICAtIGAxYDog5pyJ5pWI6K+B5Lu277yI5Liq5Lq677ya6Lqr5Lu96K+BL+aKpOeFp++8jOS8geS4mu+8muiQpeS4muaJp+eFp++8iVxuICAgKiAgIC0gYDJgOiDotKLku5jpgJrvvIjkuKrkurrvvIlcbiAgICogICAtIGAzYDog6ZO26KGM5Y2h77yI5LyB5Lia77yJXG4gICAqICAgLSBgNGA6IOW+ruS/oe+8iOS4quS6uu+8iVxuICAgKiAgIC0gYDVgOiDmiYtR77yI5Liq5Lq677yJXG4gICAqICAgLSBgNmA6IOWFrOS8l+W5s+WPsO+8iOS8geS4mu+8iVxuICAgKiAgIC0gYDdgOiDnur/kuIvorqTor4HvvIjlvojlsJHvvIlcbiAgICogICAtIGA4YDog5Zu96ZmF5L+h55So5Y2h77yI5Liq5Lq6L+S8geS4mu+8iVxuICAgKiAgIC0gYDlgOiDkvIHkuJrnur/kuIvmiZPmrL5cbiAgICogICAtIGAxMGA6IOe6v+S4iueUs+ivt++8jOe6v+S4i+WuoeaguO+8iOS8geS4muS/ruaUueWunuWQjeiupOivgea1geeoi++8iVxuICAgKiAgIC0gYDExYDog57Gz5aSn5biI6K6k6K+BXG4gICAqICAgLSBgMTJgOiDkuKrkurrkurrohLjmoLjouqvorqTor4FcbiAgICogICAtIGAyMGA6IOS7o+eQhuWVhlxuICAgKi9cbiAgYXV0aFR5cGU6IG51bWJlcjtcblxuICAvKipcbiAgICog6K6k6K+B5Zyw5Yy6XG4gICAqICAgLSBgLTFgOiDmnKrnn6VcbiAgICogICAtIGAwYDog5aSn6ZmGXG4gICAqICAgLSBgMWA6IOa4r+a+s1xuICAgKiAgIC0gYDJgOiDlj7Dmub5cbiAgICogICAtIGAzYDog5aSW57GNXG4gICAqL1xuICBhdXRoQXJlYTogbnVtYmVyO1xufVxuXG4vLyDlvZPliY3nlKjmiLfmlbDmja7vvIznvJPlrZjnrKzkuIDmrKHmi4nlj5ZcbmxldCBjdXJyZW50VXNlckRhdGE6IFByb21pc2U8QXBwVXNlckRhdGE+O1xuZXhwb3J0IGNvbnN0IGN1cnJlbnQgPSAoKSA9PlxuICBjdXJyZW50VXNlckRhdGEgfHxcbiAgKGN1cnJlbnRVc2VyRGF0YSA9IG5ldyBQcm9taXNlKHJlc29sdmUgPT5cbiAgICBfbWFuYWdlci5nZXRDb21EYXRhKChkYXRhOiBhbnkpID0+IHtcbiAgICAgIGlmICghZGF0YSkge1xuICAgICAgICByZXNvbHZlKG51bGwpO1xuICAgICAgICByZXR1cm47XG4gICAgICB9XG4gICAgICBjb25zdCB7IHVzZXJJbmZvLCBvd25lckluZm8sIGFwcElkIH0gPSBkYXRhO1xuXG4gICAgICBjb25zdCBsb2dpblVpbiA9IE51bWJlcih1c2VySW5mby51aW4pO1xuICAgICAgY29uc3Qgb3duZXJVaW4gPSBOdW1iZXIodXNlckluZm8ub3duZXJVaW4pO1xuXG4gICAgICBsZXQgaWRlbnRpdHk6IEFwcFVzZXJJZGVudGl0eUluZm8gfCBudWxsID0gbnVsbDtcblxuICAgICAgaWYgKG93bmVySW5mbyAmJiBvd25lckluZm8uYXV0aERldGFpbCAmJiBcInR5cGVcIiBpbiBvd25lckluZm8uYXV0aERldGFpbCkge1xuICAgICAgICBjb25zdCB7IHR5cGUsIGF1dGhUeXBlLCBhcmVhIH0gPSBvd25lckluZm8uYXV0aERldGFpbDtcbiAgICAgICAgLy8g5a2X5q616Kej6YeK77yaaHR0cDovL3RhcGQub2EuY29tL1FDbG91ZF8yMDE1L21hcmtkb3duX3dpa2lzLyMxMDEwMTAzOTUxMDA4MzIyNTk5XG4gICAgICAgIGlkZW50aXR5ID0ge1xuICAgICAgICAgIHN1YmplY3RUeXBlOiArdHlwZSxcbiAgICAgICAgICBhdXRoVHlwZTogK2F1dGhUeXBlLFxuICAgICAgICAgIGF1dGhBcmVhOiArYXJlYSxcbiAgICAgICAgfTtcbiAgICAgIH1cblxuICAgICAgcmVzb2x2ZSh7XG4gICAgICAgIGlzT3duZXI6IGxvZ2luVWluID09PSBvd25lclVpbixcbiAgICAgICAgbG9naW5VaW4sXG4gICAgICAgIG93bmVyVWluLFxuICAgICAgICBhcHBJZDogK2FwcElkLFxuICAgICAgICBpZGVudGl0eSxcbiAgICAgICAgbmlja05hbWU6IHVzZXJJbmZvLm5pY2ssXG4gICAgICAgIGRpc3BsYXlOYW1lOiB1c2VySW5mby5kaXNwbGF5TmFtZSxcbiAgICAgIH0pO1xuICAgIH0pXG4gICkpO1xuIl19