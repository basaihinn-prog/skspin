const fs = require("fs");
const path = require("path");

const root = process.cwd();
const now = new Date();

const runtimeIds = fs
  .readFileSync(path.join(root, "config/runtime_game_ids.txt"), "utf8")
  .split(/\r?\n/)
  .map(value => value.trim())
  .filter(Boolean);

const demoGames = JSON.parse(
  fs.readFileSync(path.join(root, "config/demoGame.json"), "utf8")
);

const gameBets = JSON.parse(
  fs.readFileSync(path.join(root, "config/GameBet.json"), "utf8")
);

const gameDb = db.getSiblingDB("game");
const adminDb = db.getSiblingDB("GameAdmin");

const demoById = new Map(
  demoGames.map(item => [String(item.id), item])
);

const betById = new Map();

for (const item of gameBets) {
  const id = String(
    item.Game_id ??
    item.GameId ??
    item.gameId ??
    item._id ??
    ""
  );

  const bet = String(
    item.Bet ??
    item.BetList ??
    item.bet ??
    ""
  );

  if (id) {
    betById.set(id, bet);
  }
}

function stripPrefix(gameId) {
  return gameId.replace(
    /^(pg|pp|jili|jdb|hacksaw|facai|spribe)_/,
    ""
  );
}

function fallbackProvider(gameId) {
  return gameId.split("_")[0].toUpperCase();
}

function findDemo(gameId) {
  return (
    demoById.get(gameId) ??
    demoById.get(stripPrefix(gameId)) ??
    null
  );
}

function findBet(gameId) {
  return (
    betById.get(gameId) ??
    betById.get(stripPrefix(gameId)) ??
    "0.1,0.2,0.3"
  );
}

const operations = runtimeIds.map(gameId => {
  const demo = findDemo(gameId);
  const betList = findBet(gameId);
  const defaultBet = Number(betList.split(",")[0]) || 0.1;

  const provider = String(
    demo?.manufacturerName ??
    fallbackProvider(gameId)
  ).toUpperCase();

  const active = demo
    ? String(demo.runStatus) === "1"
    : true;

  return {
    updateOne: {
      filter: { _id: gameId },
      update: {
        $setOnInsert: {
          _id: gameId,
          CreateAt: now
        },
        $set: {
          GameId: gameId,
          Name: gameId,
          Type: Number(demo?.gameType ?? 2),
          Status: active ? 0 : 2,
          ManufacturerName: provider,
          LineNum: 0,
          BetList: betList,
          Bet: 0,
          ChangeBetOff: 0,
          BuyType: 0,
          BuyBetMulti: 0,
          RewardPercent: 0,
          NoAwardPercent: 0,
          DefaultBet: defaultBet,
          DefaultBetLevel: gameId.startsWith("pg_") ? 1 : 0,
          GameNameConfig: {
            en: { GameName: gameId, Icon: "" },
            zh: { GameName: gameId, Icon: "" },
            th: { GameName: gameId, Icon: "" },
            idr: { GameName: gameId, Icon: "" },
            it: { GameName: gameId, Icon: "" },
            es: { GameName: gameId, Icon: "" }
          },
          Source: "runtime-source-bootstrap",
          UpdatedAt: now
        }
      },
      upsert: true
    }
  };
});

if (operations.length < 400) {
  throw new Error(
    `Runtime game catalog incomplete: ${operations.length}`
  );
}

const result = gameDb.Games.bulkWrite(
  operations,
  { ordered: true }
);

gameDb.Games.createIndex(
  { GameId: 1 },
  {
    name: "uq_games_gameid",
    unique: true
  }
);

gameDb.Games.createIndex(
  { ManufacturerName: 1, Status: 1 },
  {
    name: "ix_games_provider_status"
  }
);

adminDb.schema_migrations.updateOne(
  { _id: "006_runtime_game_master" },
  {
    $setOnInsert: {
      version: 6,
      name: "runtime_game_master",
      appliedAt: now
    },
    $set: {
      status: "applied",
      source: "gameserver_runtime_identifiers",
      gameCount: runtimeIds.length,
      updatedAt: now
    }
  },
  { upsert: true }
);

printjson({
  requested: runtimeIds.length,
  inserted: result.upsertedCount,
  matched: result.matchedCount,
  modified: result.modifiedCount,
  stored: gameDb.Games.countDocuments({})
});

print("STEP 6 RUNTIME GAME MASTER APPLIED");
