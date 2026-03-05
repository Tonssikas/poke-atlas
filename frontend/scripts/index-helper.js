/*
A helper script made for building pokemon-index.json for frontend search functionality.
The script fetches pokemon names and pokeapi urls from pokeAPI, from which we parse name and pokemon ID
*/


import fs from "fs/promises";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

async function fetchData() {
  try {
    const response = await fetch("https://pokeapi.co/api/v2/pokemon?limit=2000");
    const data = await response.json();

    const formatted = data.results.map((p) => {
      const id = p.url.split("/").filter(Boolean).pop();

      return {
        id: Number(id),
        name: p.name
      };
    });

      const filePath = path.join(__dirname, "../src/lib/data/pokemon-index.json");

    await fs.writeFile(
      filePath,
      JSON.stringify(formatted, null, 2)
    );

    console.log("Data written successfully");
  } catch (err) {
    console.error("Error:", err);
  }
}

fetchData();