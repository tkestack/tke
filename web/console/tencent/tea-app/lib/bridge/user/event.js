"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
var EventEmitter = require("eventemitter3");
var UserEmitter = /** @class */ (function (_super) {
    __extends(UserEmitter, _super);
    function UserEmitter() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    return UserEmitter;
}(EventEmitter));
exports.UserEmitter = UserEmitter;
exports.userEmitter = new UserEmitter();
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXZlbnQuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvYnJpZGdlL3VzZXIvZXZlbnQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7O0FBQUEsNENBQThDO0FBVTlDO0lBQWlDLCtCQUE0QjtJQUE3RDs7SUFBK0QsQ0FBQztJQUFELGtCQUFDO0FBQUQsQ0FBQyxBQUFoRSxDQUFpQyxZQUFZLEdBQW1CO0FBQW5ELGtDQUFXO0FBRVgsUUFBQSxXQUFXLEdBQUcsSUFBSSxXQUFXLEVBQUUsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCAqIGFzIEV2ZW50RW1pdHRlciBmcm9tIFwiZXZlbnRlbWl0dGVyM1wiO1xuXG5pbnRlcmZhY2UgSW52YWxpZGF0ZUV2ZW50QXJncyB7XG4gIHNvdXJjZTogXCJhY2NvdW50Q2hhbmdlZFwiIHwgXCJsb2dvdXRcIjtcbn1cblxuaW50ZXJmYWNlIFVzZXJFdmVudFR5cGVzIHtcbiAgaW52YWxpZGF0ZTogW0ludmFsaWRhdGVFdmVudEFyZ3NdO1xufVxuXG5leHBvcnQgY2xhc3MgVXNlckVtaXR0ZXIgZXh0ZW5kcyBFdmVudEVtaXR0ZXI8VXNlckV2ZW50VHlwZXM+IHt9XG5cbmV4cG9ydCBjb25zdCB1c2VyRW1pdHRlciA9IG5ldyBVc2VyRW1pdHRlcigpO1xuIl19