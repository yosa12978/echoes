// get all card elements.
const cards = document.querySelectorAll(".card");

// create colorthief instance
const colorThief = new ColorThief();

cards.forEach(async (card) => {
    const image = card.children[0];
    const text = card.children[1];
    // get palette color from image
    const palette = await extractColor(image);

    const primary = palette[0].join(",");
    const secondary = palette[1].join(",");

    // change color
    card.style.background = `rgb(${primary})`;
    text.style.color = `rgb(${secondary})`;
});

// async function wrapper
function extractColor(image) {
    return new Promise((resolve) => {
        const getPalette = () => {
            return colorThief.getPalette(image, 4);
        };

        // as said in the colorthief documentation, 
        // we have to wait until the image is fully loaded.

        if (image.complete) {
            return resolve(getPalette());
        }

        image.onload = () => {
            resolve(getPalette());
        };
    });
}
