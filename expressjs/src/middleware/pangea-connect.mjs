import { PangeaConfig, IPIntelService, PangeaErrors } from "pangea-node-sdk";
import redis from 'redis';
const domain = process.env.PANGEA_DOMAIN;
const token = process.env.PANGEA_INTEL_TOKEN;
const config = new PangeaConfig({ domain: domain });
const ipIntel = new IPIntelService(String(token), config);
const ipIntelType = process.env.PANGEA_IP_INTEL_TYPE;
const ipIntelThreshold = process.env.PANGEA_IP_INTEL_SCORE_THRESHOLD;
const redisUrl = process.env.REDIS_URL;
const redisPassword = process.env.REDIS_PASSWORD;
const redisTtl = process.env.REDIS_TTL;
const options = { provider: "crowdstrike", verbose: true, raw: true };

const redisClient = await redis.createClient({
    url: 'redis://' + redisUrl,
    connectTimeout: 5000,
}).on('error', err => console.log('Redis Client Error', err))
    .connect();

export const PangeaConnect = (req, res, next) => {

    if (!ipIntelType || !token || !domain) {
        next();
        return;
    }
    let ipAddress = "";
    if (ipIntelType === "header") {
        ipAddress = req.headers['x-forwarded-for'] || req.connection.remoteAddress;
    } else if (ipIntelType === "request") {
        ipAddress = req.connection.remoteAddress;
    }

    redisClient.get(ipAddress).then((isBlocked) => {
        allowOrDenyRequest(req, res, next, isBlocked, ipAddress);
    }).catch((err) => {
        next();
        return;
    });
}

const allowOrDenyRequest = (req, res, next, isBlocked, ipAddress) => {
    if (isBlocked === null) {
        processIpIntel(ipAddress);
        next();
    } else if (isBlocked === "Y") {
        res.sendStatus(403);
    } else {
        next();
    }
};

const processIpIntel = async (ipAddress) => {
    try {
        const response = await ipIntel.reputation(ipAddress, options);
        if (response.status.toLowerCase() === "success") {
            let blockedStatus = response.result.data.score >= ipIntelThreshold ? "Y" : "N";
            redisClient.set(ipAddress, blockedStatus, {
                EX: redisTtl
            });
        } else {
            redisClient.set(ipAddress, "N", {
                EX: redisTtl
            });

        }
    } catch (ignored) { }

}