
function main() {
    let toggleButton = document.getElementById("portrait-nav-toggle")

    toggleButton.addEventListener("click", function () {
        portraitToggleHandler(toggleButton)
    })
}

function portraitToggleHandler(event) {
    window.scrollTo(0, 0)

    const classDisplayNone = "my-display-none"
    let portraitNavMenuElement = document.getElementById("my-portrait-nav-menu")

    let displayingPortraitMenu = portraitNavMenuElement.classList.contains(classDisplayNone)
    document.body.style.overflow = displayingPortraitMenu ? "hidden" : ""

    if (displayingPortraitMenu) {
        portraitNavMenuElement.classList.remove(classDisplayNone)
        return
    }

    portraitNavMenuElement.classList.add(classDisplayNone);
}

window.onload = main