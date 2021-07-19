"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var FrequencyLimiter = /** @class */ (function () {
    function FrequencyLimiter(_a) {
        var _b = _a.name, name = _b === void 0 ? "" : _b, _c = _a.limitTimes, limitTimes = _c === void 0 ? 10 : _c, _d = _a.tenseLimitTimes, tenseLimitTimes = _d === void 0 ? 100 : _d, _e = _a.timeout, timeout = _e === void 0 ? 10 * 1000 : _e;
        this.times = 0;
        this.name = name;
        this.timeout = timeout;
        this.limitTimes = limitTimes;
        this.tenseLimitTimes = tenseLimitTimes;
    }
    FrequencyLimiter.prototype.exec = function (fn) {
        var _this = this;
        if (this.tense) {
            return;
        }
        if (!this.timer) {
            this.timer = window.setTimeout(function () {
                _this.times = 0;
                _this.timer = null;
            }, this.timeout);
        }
        this.times++;
        if (this.times <= this.limitTimes) {
            fn();
        }
        else if (!this.hasOutput) {
            this.hasOutput = true;
            console.error(this.name + " \u8C03\u7528\u9891\u7387\u8D85\u8FC7\u9650\u5236");
        }
        if (this.times > this.tenseLimitTimes) {
            this.tense = true;
        }
    };
    return FrequencyLimiter;
}());
exports.FrequencyLimiter = FrequencyLimiter;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZnJlcXVlbmN5TGltaXRlci5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL3NyYy9oZWxwZXJzL2ZyZXF1ZW5jeUxpbWl0ZXIudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQTtJQW1CRSwwQkFBWSxFQUtYO1lBSkMsWUFBUyxFQUFULDhCQUFTLEVBQ1Qsa0JBQWUsRUFBZixvQ0FBZSxFQUNmLHVCQUFxQixFQUFyQiwwQ0FBcUIsRUFDckIsZUFBbUIsRUFBbkIsd0NBQW1CO1FBRW5CLElBQUksQ0FBQyxLQUFLLEdBQUcsQ0FBQyxDQUFDO1FBQ2YsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUM7UUFDakIsSUFBSSxDQUFDLE9BQU8sR0FBRyxPQUFPLENBQUM7UUFDdkIsSUFBSSxDQUFDLFVBQVUsR0FBRyxVQUFVLENBQUM7UUFDN0IsSUFBSSxDQUFDLGVBQWUsR0FBRyxlQUFlLENBQUM7SUFDekMsQ0FBQztJQUVNLCtCQUFJLEdBQVgsVUFBWSxFQUFZO1FBQXhCLGlCQXVCQztRQXRCQyxJQUFJLElBQUksQ0FBQyxLQUFLLEVBQUU7WUFDZCxPQUFPO1NBQ1I7UUFFRCxJQUFJLENBQUMsSUFBSSxDQUFDLEtBQUssRUFBRTtZQUNmLElBQUksQ0FBQyxLQUFLLEdBQUcsTUFBTSxDQUFDLFVBQVUsQ0FBQztnQkFDN0IsS0FBSSxDQUFDLEtBQUssR0FBRyxDQUFDLENBQUM7Z0JBQ2YsS0FBSSxDQUFDLEtBQUssR0FBRyxJQUFJLENBQUM7WUFDcEIsQ0FBQyxFQUFFLElBQUksQ0FBQyxPQUFPLENBQUMsQ0FBQztTQUNsQjtRQUVELElBQUksQ0FBQyxLQUFLLEVBQUUsQ0FBQztRQUNiLElBQUksSUFBSSxDQUFDLEtBQUssSUFBSSxJQUFJLENBQUMsVUFBVSxFQUFFO1lBQ2pDLEVBQUUsRUFBRSxDQUFDO1NBQ047YUFBTSxJQUFJLENBQUMsSUFBSSxDQUFDLFNBQVMsRUFBRTtZQUMxQixJQUFJLENBQUMsU0FBUyxHQUFHLElBQUksQ0FBQztZQUN0QixPQUFPLENBQUMsS0FBSyxDQUFJLElBQUksQ0FBQyxJQUFJLHNEQUFXLENBQUMsQ0FBQztTQUN4QztRQUVELElBQUksSUFBSSxDQUFDLEtBQUssR0FBRyxJQUFJLENBQUMsZUFBZSxFQUFFO1lBQ3JDLElBQUksQ0FBQyxLQUFLLEdBQUcsSUFBSSxDQUFDO1NBQ25CO0lBQ0gsQ0FBQztJQUNILHVCQUFDO0FBQUQsQ0FBQyxBQXhERCxJQXdEQztBQXhEWSw0Q0FBZ0IiLCJzb3VyY2VzQ29udGVudCI6WyJleHBvcnQgY2xhc3MgRnJlcXVlbmN5TGltaXRlciB7XG4gIHByaXZhdGUgbmFtZTogc3RyaW5nO1xuICBwcml2YXRlIHRpbWVyOiBudW1iZXI7XG5cbiAgLy8g5b2T5YmN6LCD55So5qyh5pWwXG4gIHByaXZhdGUgdGltZXM6IG51bWJlcjtcblxuICAvLyDml7bpl7TpmZDml7ZcbiAgcHJpdmF0ZSB0aW1lb3V0OiBudW1iZXI7XG4gIC8vIOasoeaVsOmZkOWItlxuICBwcml2YXRlIGxpbWl0VGltZXM6IG51bWJlcjtcbiAgLy8g5pyA6auY5qyh5pWw6ZmQ5Yi277yM6LaF5Ye65bCG5ouS57ud5ZCO57ut5omA5pyJ5omn6KGMXG4gIHByaXZhdGUgdGVuc2VMaW1pdFRpbWVzOiBudW1iZXI7XG5cbiAgLy8g5piv5ZCm5bey57uP6LaF5Ye65pyA6auY6ZmQ5Yi2XG4gIHByaXZhdGUgdGVuc2U6IGJvb2xlYW47XG4gIC8vIOaYr+WQpuW3sue7j+i+k+WHuuWRiuitplxuICBwcml2YXRlIGhhc091dHB1dDogYm9vbGVhbjtcblxuICBjb25zdHJ1Y3Rvcih7XG4gICAgbmFtZSA9IFwiXCIsXG4gICAgbGltaXRUaW1lcyA9IDEwLFxuICAgIHRlbnNlTGltaXRUaW1lcyA9IDEwMCxcbiAgICB0aW1lb3V0ID0gMTAgKiAxMDAwLFxuICB9KSB7XG4gICAgdGhpcy50aW1lcyA9IDA7XG4gICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICB0aGlzLnRpbWVvdXQgPSB0aW1lb3V0O1xuICAgIHRoaXMubGltaXRUaW1lcyA9IGxpbWl0VGltZXM7XG4gICAgdGhpcy50ZW5zZUxpbWl0VGltZXMgPSB0ZW5zZUxpbWl0VGltZXM7XG4gIH1cblxuICBwdWJsaWMgZXhlYyhmbjogRnVuY3Rpb24pIHtcbiAgICBpZiAodGhpcy50ZW5zZSkge1xuICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghdGhpcy50aW1lcikge1xuICAgICAgdGhpcy50aW1lciA9IHdpbmRvdy5zZXRUaW1lb3V0KCgpID0+IHtcbiAgICAgICAgdGhpcy50aW1lcyA9IDA7XG4gICAgICAgIHRoaXMudGltZXIgPSBudWxsO1xuICAgICAgfSwgdGhpcy50aW1lb3V0KTtcbiAgICB9XG5cbiAgICB0aGlzLnRpbWVzKys7XG4gICAgaWYgKHRoaXMudGltZXMgPD0gdGhpcy5saW1pdFRpbWVzKSB7XG4gICAgICBmbigpO1xuICAgIH0gZWxzZSBpZiAoIXRoaXMuaGFzT3V0cHV0KSB7XG4gICAgICB0aGlzLmhhc091dHB1dCA9IHRydWU7XG4gICAgICBjb25zb2xlLmVycm9yKGAke3RoaXMubmFtZX0g6LCD55So6aKR546H6LaF6L+H6ZmQ5Yi2YCk7XG4gICAgfVxuXG4gICAgaWYgKHRoaXMudGltZXMgPiB0aGlzLnRlbnNlTGltaXRUaW1lcykge1xuICAgICAgdGhpcy50ZW5zZSA9IHRydWU7XG4gICAgfVxuICB9XG59XG4iXX0=