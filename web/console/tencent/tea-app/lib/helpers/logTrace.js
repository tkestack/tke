"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * 打印异常信息到浏览器控制台
 */
function logTrace(leadMessage, trace) {
    if (console.groupCollapsed) {
        console.groupCollapsed("%c%s", "background-color: #fdf0f0; color: red;", "" + leadMessage);
    }
    else {
        console.log("%c%s", "background-color: #fdf0f0; color: red;", "" + leadMessage);
    }
    console.log(trace.stack);
    if (console.groupEnd) {
        console.groupEnd();
    }
}
exports.logTrace = logTrace;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibG9nVHJhY2UuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvaGVscGVycy9sb2dUcmFjZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBOztHQUVHO0FBQ0gsU0FBZ0IsUUFBUSxDQUFDLFdBQW1CLEVBQUUsS0FBSztJQUNqRCxJQUFJLE9BQU8sQ0FBQyxjQUFjLEVBQUU7UUFDMUIsT0FBTyxDQUFDLGNBQWMsQ0FDcEIsTUFBTSxFQUNOLHdDQUF3QyxFQUN4QyxLQUFHLFdBQWEsQ0FDakIsQ0FBQztLQUNIO1NBQU07UUFDTCxPQUFPLENBQUMsR0FBRyxDQUNULE1BQU0sRUFDTix3Q0FBd0MsRUFDeEMsS0FBRyxXQUFhLENBQ2pCLENBQUM7S0FDSDtJQUNELE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLEtBQUssQ0FBQyxDQUFDO0lBQ3pCLElBQUksT0FBTyxDQUFDLFFBQVEsRUFBRTtRQUNwQixPQUFPLENBQUMsUUFBUSxFQUFFLENBQUM7S0FDcEI7QUFDSCxDQUFDO0FBbEJELDRCQWtCQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICog5omT5Y2w5byC5bi45L+h5oGv5Yiw5rWP6KeI5Zmo5o6n5Yi25Y+wXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBsb2dUcmFjZShsZWFkTWVzc2FnZTogc3RyaW5nLCB0cmFjZSkge1xuICBpZiAoY29uc29sZS5ncm91cENvbGxhcHNlZCkge1xuICAgIGNvbnNvbGUuZ3JvdXBDb2xsYXBzZWQoXG4gICAgICBcIiVjJXNcIixcbiAgICAgIFwiYmFja2dyb3VuZC1jb2xvcjogI2ZkZjBmMDsgY29sb3I6IHJlZDtcIixcbiAgICAgIGAke2xlYWRNZXNzYWdlfWBcbiAgICApO1xuICB9IGVsc2Uge1xuICAgIGNvbnNvbGUubG9nKFxuICAgICAgXCIlYyVzXCIsXG4gICAgICBcImJhY2tncm91bmQtY29sb3I6ICNmZGYwZjA7IGNvbG9yOiByZWQ7XCIsXG4gICAgICBgJHtsZWFkTWVzc2FnZX1gXG4gICAgKTtcbiAgfVxuICBjb25zb2xlLmxvZyh0cmFjZS5zdGFjayk7XG4gIGlmIChjb25zb2xlLmdyb3VwRW5kKSB7XG4gICAgY29uc29sZS5ncm91cEVuZCgpO1xuICB9XG59XG4iXX0=