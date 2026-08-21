export const GeneratorDefaultCs = (minBet, maxBet):any[] => {



    // let curve = [0.4];  // 预定义的曲线数组，用于调整金币的变化
    let curve = [.05, .1, .2, .4];  // 预定义的曲线数组，用于调整金币的变化
    let levels = 10

    while (minBet * levels < maxBet / levels * curve[0]){

        curve.unshift(curve[0] * 0.2);
    }


    // 如果最大投注金额与最小投注金额的比例小于级别数，则重新计算级别数
    if (maxBet / minBet < levels) {
        levels = maxBet * 1E3 / (minBet * 1E3) | 0;
    }

    // 计算最大金币值
    let maxCoinValue = maxBet * 1E3 / levels / 1E3;

    let d = maxBet * 1E3 / levels / 1E3;
    while (Math.floor((maxCoinValue + 1E-4) * 100) / 100 * levels < maxBet) {
        levels--;
        maxCoinValue = maxBet * 1E3 / levels / 1E3;
    }


    // 如果需要进一步处理金币值，调用 CoinManager 的方法进行处理
    maxCoinValue = getNiceCoinValue(maxCoinValue)


    // 将最大金币值取整到两位小数
    maxCoinValue = Math.floor((maxCoinValue + 1E-4) * 100) / 100;

    // 确保最大金币值在允许范围内
    if (maxCoinValue * levels > maxBet)
        maxCoinValue = (maxCoinValue * 100 | 0) / 100;

    // 计算金币值范围
    var x = maxCoinValue - minBet;
    var coinValues = [];
    coinValues.push(minBet / 10000);

    // 根据曲线数组计算并添加金币值
    for (var j = 0; j < curve.length; j++) {

        // let a = minBet + x * curve[j]
        var computedVal = getNiceCoinValue(minBet + x * curve[j]);
        if (computedVal > minBet && computedVal < maxCoinValue){

            let num = (computedVal / 10000).toFixed(2)
            coinValues.push(getNiceCoinValue(parseFloat(num)));
        }


    }
    // 添加最大金币值
    coinValues.push(d / 10000);

    // 删除相近的重复值
    for (var i = 1; i < coinValues.length; i++)
        if (Math.abs(coinValues[i] - coinValues[i - 1]) < .001) {
            coinValues.splice(i, 1);
            i--;
        }

    // 计算最大有效投注金额，基于金币的最大值、投注级别（levels）和最大线路数（linesForMaxBet）

    let maxValidBetC = -1;  // 初始化最大有效投注金额为 -1，后续将通过循环更新
    // 遍历每个投注级别和每个金币值
    for (var lvl = 1; lvl <= levels; lvl++)  // 循环所有投注级别
        for (var c = 0; c < coinValues.length; c++) {  // 循环所有金币值
            // 计算在当前级别和最小线路数情况下的投注金额
            var betWithMinLinesC = Math.round(coinValues[c] * 100 * lvl * minBet);
            // 计算在当前级别和最大线路数情况下的投注金额
            var betWithMaxLinesC = Math.round(coinValues[c] * 100 * lvl * maxBet);
            // 如果最小线路投注金额不超过最大投注限制，并且比当前最大有效投注金额大，则更新最大有效投注金额
            if (betWithMinLinesC <= maxBet)
                if (betWithMinLinesC > maxValidBetC) maxValidBetC = betWithMinLinesC;
            // 如果最大线路投注金额不超过最大投注限制，并且比当前最大有效投注金额大，则更新最大有效投注金额
            if (betWithMaxLinesC <= maxBet)
                if (betWithMaxLinesC > maxValidBetC) maxValidBetC = betWithMaxLinesC;
        }



    return coinValues
}

const getNiceCoinValue = (value) => {
    value = parseFloat(parseFloat(value).toFixed(2))

    var niceValue;
    // if (ServerOptions.amountType == "COINS_V2") {
    //     var log10val = Math.floor(Math.log10(value));
    //     var pow10 = Math.pow(10, log10val | 0);
    //     var niceSubUnitSteps = [1, 2, 5];
    //     var niceSubUnitValues = [];
    //     for (var i = 0; i < niceSubUnitSteps.length; i++) niceSubUnitValues.push(pow10 * niceSubUnitSteps[i]);
    //     var minDist = Math.pow(10, log10val + 1);
    //     var index = -1;
    //     for (var i = 0; i < niceSubUnitValues.length; i++) {
    //         var d = Math.abs(value - niceSubUnitValues[i]);
    //         if (d < minDist) {
    //             minDist = d;
    //             index = i
    //         }
    //     }
    //     niceValue = niceSubUnitValues[index]
    // }
    if (value > 5)
        niceValue = Math.floor(value);
    else if (value >= 1)
        niceValue = Math.floor(value * 4) / 4;
    else {
        var niceSubUnitValues = [.01, .02, .03, .05, .07, .1, .2, .3, .5, .75];
        var minDist = 5;
        var index = -1;
        for (var i = 0; i < niceSubUnitValues.length; i++) {
            var d = Math.abs(value - niceSubUnitValues[i]);
            if (d < minDist) {
                minDist = d;
                index = i
            }
        }
        niceValue = niceSubUnitValues[index]
    }
    return niceValue
};
