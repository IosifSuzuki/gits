"use strict";

(function () {

    function degreeToRadian(degree) {
        return degree * (2 * Math.PI / 360)
    }

    function radianToDegree(radian) {
        return radian * (360 / (2 * Math.PI))
    }

    // We would work with math circle
    function normalizeAngle(angle) {
        return ((angle - 90) + 360) % 360
    }

    function polarToCartesian(dist, azim) {
        let phi = degreeToRadian(normalizeAngle(azim))

        return {
            x: dist * Math.cos(phi),
            y: dist * Math.sin(phi)
        }
    }

    function distance(ax, ay, bx, by) {
        let xDistance = Math.abs(ax - bx)
        let yDistance = Math.abs(ay - by)

        return Math.sqrt(xDistance * xDistance + yDistance * yDistance)
    }

    function degreeAlphaFromLawOfCosines(a, b, c) {
        let radian = Math.acos((b * b + c * c - a * a) / (2 * b * c))

        return radianToDegree(radian)
    }

    function getYFromLine(x2, y2, x1, y1, x) {
        return (x - x1) / (x2 - x1) * (y2 - y1) + y1
    }

    function signOffsetDegree(aCoord, bCoord) {
        let y = getYFromLine(aCoord.x, aCoord.y, 0, 0, bCoord.x)

        let isNegativeX = aCoord.x < 0
        if (isNegativeX) {
            return y > bCoord.y ? -1 : 1
        }
        return y > bCoord.y ? 1 : -1
    }

    function displayResult(targetPolarCoord, correctedPolarCoord) {
        displayHtmlElement(messageResultElement, true)

        let message = "Distance to target: " + targetPolarCoord.r.toFixed(2) + "; Azimuth to target: " + targetPolarCoord.azimuth.toFixed(2) + "<br />"
        if (correctedPolarCoord) {
            message += "Corrected distance to target: " + correctedPolarCoord.r.toFixed(2) + "; Corrected azimuth to target: " + correctedPolarCoord.azimuth.toFixed(2)
        }

        printMessage(message)
    }

    function printMessage(text) {
        messageResultElement.innerHTML = text
    }

    function hasHiddenHtmlElement(element) {
        return element.classList.contains(classDisplayNone)
    }

    function displayHtmlElement(element, display) {
        if (hasHiddenHtmlElement(element) && display) {
            element.classList.remove(classDisplayNone)
        } else if (!hasHiddenHtmlElement(element) && !display) {
            element.classList.add(classDisplayNone)
        }
    }

    function applyTranslate(aCoord, bCoord) {
        aCoord.x += bCoord.x;
        aCoord.y += bCoord.y;
    }

    function calculateTargetLocation(enemyDistance, enemyAzimuth, artilleryDistance, artilleryAzimuth) {
        let enemy = polarToCartesian(enemyDistance, enemyAzimuth)
        let artillery = polarToCartesian(artilleryDistance, artilleryAzimuth)

        let r = distance(enemy.x, enemy.y, artillery.x, artillery.y)
        let deltaDegree = degreeAlphaFromLawOfCosines(enemyDistance, artilleryDistance, r)

        let sign = signOffsetDegree(artillery, enemy)
        let oppositeArtilleryToSpotter = (artilleryAzimuth + 180) % 360

        let azimuth = (oppositeArtilleryToSpotter + sign * deltaDegree) % 360

        if (azimuth < 0) {
            azimuth = 360 + azimuth
        }

        return {
            r: r,
            azimuth: azimuth,
        }
    }

    function calculateWindCoorection(enemyPolarCoord, windLevel, windAzimuth, artilleryType) {
        const offsetArtillery = {
            1: 25,
            2: 25,
            3: 50,
        }

        let windR = offsetArtillery[artilleryType] * windLevel
        let enemy = polarToCartesian(enemyPolarCoord.r, enemyPolarCoord.azimuth)

        let wind = polarToCartesian(windR, windAzimuth)
        applyTranslate(wind, enemy)

        let correctionR = distance(0, 0, wind.x, wind.y)
        let deltaDegree = degreeAlphaFromLawOfCosines(windR, correctionR, enemyPolarCoord.r)
        let sign = signOffsetDegree(enemy, wind)

        let azimuth = (enemyPolarCoord.azimuth + sign * deltaDegree) % 360
        if (azimuth < 0) {
            azimuth = 360 + azimuth
        }

        return {
            r: correctionR,
            azimuth: azimuth,
        }
    }

    const classDisplayNone = "d-none"

    let messageResultElement = document.getElementById("message-result")
    let instructionElement = document.getElementById("instruction")
    let computeButton = document.getElementById("compute-button")
    let clearButton = document.getElementById("clear-button")
    let windDirectionButton = document.getElementById("wind-direction-button")
    let windCorrectionSection = document.getElementById("wind-correction-section")
    let enemyDistanceInput = document.getElementById("target_distance")
    let enemyAzimuthInput = document.getElementById("target_azimuth")
    let artilleryDistanceInput = document.getElementById("artillery_distance")
    let artilleryAzimuthInput = document.getElementById("artillery_azimuth")
    let windLevelInput = document.getElementById("wind_level")
    let windAzimuthInput = document.getElementById("wind_azimuth")
    let instructionButton = document.getElementById("instruction-button")
    let artilleryTypeInput = document.getElementById("artillery-type")

    instructionButton.addEventListener("click", function() {
        displayHtmlElement(instructionElement, hasHiddenHtmlElement(instructionElement))
    })

    windDirectionButton.addEventListener("click", function() {
        displayHtmlElement(windCorrectionSection, hasHiddenHtmlElement(windCorrectionSection))
    })

    clearButton.addEventListener("click", function() {
        enemyDistanceInput.value = ""
        enemyAzimuthInput.value = ""
        artilleryDistanceInput.value = ""
        artilleryAzimuthInput.value = ""
        windLevelInput.value = "0"
        windAzimuthInput.value = ""

        displayHtmlElement(windCorrectionSection, false)
    })

    computeButton.addEventListener("click", function() {
        let enemyDistance = Math.round(enemyDistanceInput.value)
        let enemyAzimuth = Math.round(enemyAzimuthInput.value)
        let artilleryDistance = Math.round(artilleryDistanceInput.value)
        let artilleryAzimuth = Math.round(artilleryAzimuthInput.value)

        let targetPolarCoord = calculateTargetLocation(enemyDistance, enemyAzimuth, artilleryDistance, artilleryAzimuth)

        let correctedPolarCoord = null
        if (!hasHiddenHtmlElement(windCorrectionSection)) {
            let windLevel = Math.round(windLevelInput.value)
            let windAzimuth = Math.round(windAzimuthInput.value)
            let artilleryType = Math.round(artilleryTypeInput.value)

            correctedPolarCoord = calculateWindCoorection(targetPolarCoord, windLevel, windAzimuth, artilleryType)
        }

        displayResult(targetPolarCoord, correctedPolarCoord)
    })
})()