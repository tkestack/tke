"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports._bridge = function (moduleName) {
    return window.seajs && window.seajs.require(moduleName);
};
exports._manager = exports._bridge("manager");
exports._util = exports._bridge("util");
exports._appUtil = exports._bridge("appUtil");
exports._sdk = exports._bridge("nmc/sdk/sdk");
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiX2JyaWRnZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL3NyYy9icmlkZ2UvX2JyaWRnZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFhLFFBQUEsT0FBTyxHQUFHLFVBQUMsVUFBa0I7SUFDeEMsT0FBQSxNQUFNLENBQUMsS0FBSyxJQUFJLE1BQU0sQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLFVBQVUsQ0FBQztBQUFoRCxDQUFnRCxDQUFDO0FBQ3RDLFFBQUEsUUFBUSxHQUFHLGVBQU8sQ0FBQyxTQUFTLENBQUMsQ0FBQztBQUM5QixRQUFBLEtBQUssR0FBRyxlQUFPLENBQUMsTUFBTSxDQUFDLENBQUM7QUFDeEIsUUFBQSxRQUFRLEdBQUcsZUFBTyxDQUFDLFNBQVMsQ0FBQyxDQUFDO0FBQzlCLFFBQUEsSUFBSSxHQUFHLGVBQU8sQ0FBQyxhQUFhLENBQUMsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImV4cG9ydCBjb25zdCBfYnJpZGdlID0gKG1vZHVsZU5hbWU6IHN0cmluZykgPT5cbiAgd2luZG93LnNlYWpzICYmIHdpbmRvdy5zZWFqcy5yZXF1aXJlKG1vZHVsZU5hbWUpO1xuZXhwb3J0IGNvbnN0IF9tYW5hZ2VyID0gX2JyaWRnZShcIm1hbmFnZXJcIik7XG5leHBvcnQgY29uc3QgX3V0aWwgPSBfYnJpZGdlKFwidXRpbFwiKTtcbmV4cG9ydCBjb25zdCBfYXBwVXRpbCA9IF9icmlkZ2UoXCJhcHBVdGlsXCIpO1xuZXhwb3J0IGNvbnN0IF9zZGsgPSBfYnJpZGdlKFwibm1jL3Nkay9zZGtcIik7XG4iXX0=