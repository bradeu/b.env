const crypto = require("crypto");

const algo = "aes-256-cbc"
const secretKey = "benbenbenbenben11234567890abcdef" // FOR MVP PURPOSES, OTHERWISE WE WOULD USE CLOUD TO STORE THE SECRET AUTH KEY

const iv = crypto.randomBytes(16);

function encrypt(data) {
    const cipher = crypto.createCipheriv(algo, Buffer.from(secretKey), iv);
    let encrypted = cipher.update(data, "utf8", "hex");
    encrypted += cipher.final("hex");
    return encrypted;  
  }

function decrypt(data) {
    const decipher = crypto.createDecipheriv(algo, Buffer.from(secretKey), iv);
    let decrypted = decipher.update(data, "utf8", "hex");
    decrypted += decipher.final("hex");
    return decrypted;  
}

export default {encrypt, decrypt}