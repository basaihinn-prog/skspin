const adminDb = db.getSiblingDB("GameAdmin");
const gameDb = db.getSiblingDB("game");
const now = new Date();
const appId = "asiagaming";

const operator = adminDb.AdminOperator.findOne({ AppID: appId });

if (!operator) {
  throw new Error("AsiaGaming operator is missing");
}

adminDb.AdminOperator.updateOne(
  { AppID: appId },
  {
    $set: {
      OperatorType: 2,
      Currency: "MMK",
      WalletMode: 1,
      Status: 0,
      UpdatedAt: now
    }
  }
);

const games = gameDb.Games.find({}).toArray();

if (games.length !== 413) {
  throw new Error(`Expected 413 games, found ${games.length}`);
}

const operations = games.map(game => {
  const betBase = String(
    game.BetList || "0.1,0.2,0.3"
  );

  const bets = betBase
    .split(",")
    .map(value => Number(value))
    .filter(value => Number.isFinite(value));

  const minimumBet = bets[0] || 0.1;
  const maximumBet = bets[bets.length - 1] || minimumBet;

  const config = {
    AppID: appId,
    AppName: operator.Name || "AsiaGaming",
    GameId: game._id,
    GameName: game.Name || game._id,
    ConfigPath: "config/path",
    StopLoss: 1,
    MaxMultipleOff: 0,
    MaxMultiple: 100,
    MaxWinPoints: 1000000,
    BetBase: betBase,
    GamePattern: 3,
    Preset: 10,
    GameOn: Number(game.Status) === 0 ? 0 : 2,
    ShowNameAndTimeOff: 1,
    ShowExitBtnOff: 0,
    ExitLink: "",
    RTP: 93,
    BuyRTP: 90,
    BuyMinAwardPercent: 0,
    OnlineUpNum: minimumBet,
    OnlineDownNum: maximumBet,
    ProfitMargin: 3,
    CrashRate: 1,
    Scale: 3,
    RewardPercent: 20,
    NoAwardPercent: 200,
    UpdatedAt: now
  };

  if (String(game._id).startsWith("pg_")) {
    config.DefaultCs =
      Number(game.DefaultBet) || minimumBet;
    config.DefaultBetLevel =
      Number(game.DefaultBetLevel) || 1;
  }

  if (String(game._id).startsWith("hacksaw_")) {
    config.DefaultCs =
      Number(game.DefaultBet) || minimumBet;
  }

  return {
    updateOne: {
      filter: {
        AppID: appId,
        GameId: game._id
      },
      update: {
        $setOnInsert: {
          CreatedAt: now
        },
        $set: config
      },
      upsert: true
    }
  };
});

const result = adminDb.GameConfig.bulkWrite(
  operations,
  { ordered: true }
);

adminDb.schema_migrations.updateOne(
  { _id: "007_asiagaming_game_config" },
  {
    $setOnInsert: {
      version: 7,
      name: "asiagaming_game_config",
      appliedAt: now
    },
    $set: {
      status: "applied",
      AppID: appId,
      gameCount: games.length,
      updatedAt: now
    }
  },
  { upsert: true }
);

const stored = adminDb.GameConfig.countDocuments({
  AppID: appId
});

const normalizedOperator =
  adminDb.AdminOperator.findOne({ AppID: appId });

if (stored !== games.length) {
  throw new Error(
    `GameConfig count mismatch: ${stored}/${games.length}`
  );
}

if (Number(normalizedOperator.OperatorType) !== 2) {
  throw new Error("OperatorType normalization failed");
}

printjson({
  AppID: appId,
  requested: games.length,
  inserted: result.upsertedCount,
  matched: result.matchedCount,
  modified: result.modifiedCount,
  stored,
  operatorType: normalizedOperator.OperatorType
});

print("STEP 7 ASIAGAMING GAME CONFIG APPLIED");
