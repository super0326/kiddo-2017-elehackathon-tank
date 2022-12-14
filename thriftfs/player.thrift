
/**
 * Types and Structures
 */

struct Position {
    1:i32 x,
    2:i32 y
}

enum Direction {
    UP = 1,
    DOWN = 2,
    LEFT = 3,
    RIGHT = 4
}

struct Tank {
    1:i32 id,
    2:Position pos,
    3:Direction dir,
    4:i32 hp
}

struct Shell {
    1:i32 id,
    2:Position pos,
    3:Direction dir,
}

struct GameState {
    1:list<Tank> tanks,
    2:list<Shell> shells,
    3:i32 yourFlagNo,
    4:i32 enemyFlagNo,
    5: optional Position flagPos
}

struct Order {
    /**
    * DO NOT try to send a order with competitor's tank id.
    * In that case, game engine will treat it as cheat and would ignore ALL this player's orders in this round.
    **/
    1:i32 tankId,
    /**
    * Possible orders are: turnTo, fire, move. All others words are illegal and will be ignored.
    * If want a tank to stick around, just do NOT send any order with that tank.
    **/
    2:string order,
    /**
    * the dir are always on base of the map instead of the tank itself,
    * which mean if a 'fire' order with UP direction will made the tank fire a shell toward the UP diction of the map.
    *
    * Only move order does not need a direction, in that case just give a direction and game engine will ignore it.
    **/
    3:Direction dir
}

struct Args {
    1:i32 tankSpeed,
    2:i32 shellSpeed,
    3:i32 tankHP,
    4:i32 tankScore,
    5:i32 flagScore,
    6:i32 maxRound,
    7:i32 roundTimeoutInMs
}


/**
 * Exceptions
 */
enum PlayerErrorCode {
    UNKNOWN_ERROR = 0,
    DATABASE_ERROR = 1,
    TOO_BUSY_ERROR = 2,
}

exception PlayerUserException {
   1: required PlayerErrorCode error_code,
   2: required string error_name,
   3: optional string message,
}

exception PlayerSystemException {
   1: required PlayerErrorCode error_code,
   2: required string error_name,
   3: optional string message,
}

exception PlayerUnknownException {
   1: required PlayerErrorCode error_code,
   2: required string error_name,
   3: required string message,
}

/**
 * API
 */
service PlayerService {
    bool ping()
        throws (1: PlayerUserException user_exception,
                2: PlayerSystemException system_exception,
                3: PlayerUnknownException unknown_exception,)


    /**
    * Upload the map to player.
    * The map is made of two-dimesional array of integer. The first dimension means row of the map. The second dimension means column of the map.
    *
    * For example, if N is the map size, position(0,0) means upper left corner, position(0,N) means the upper right corner.
    * In the map array, 0 means empty field, 1 means barrier, 2 means woods, 3 means flag.
    **/
    void uploadMap(1:list<list<i32>> gamemap);


    void uploadParamters(1:Args arguments);


    /**
    * Assign a list of tank id to the player.
    * each player may have more than one tank, so the parameter is a list.
    **/
    void assignTanks(1:list<i32> tanks);


    /**
    * Report latest game state to player.
    **/
    void latestState(1:GameState state);


    /**
    * Ask for the tank orders for this round.
    * If this funtion does not return orders within the given round timeout, game engine will make all this player's tank to stick around.
    */
    list<Order> getNewOrders();
}
