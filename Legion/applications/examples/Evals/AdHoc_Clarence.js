const fs = require('fs');
const path = require('path');

// Utility to format logs
function formatLog(id, val, optName, upValue, lamportTime, isRemote = false, rid = null) {
    let logEntry = `R_ID:${id} St:${val} LT:${lamportTime} Opt:${optName}_${upValue} `;
    if (isRemote) {
        logEntry += `Type:{Remote, rid:${rid} local_LT:${rid2LT[rid].localEvent()}}\n`;
    } else {
        logEntry += `Type:Local\n`;
    }
    return logEntry;
}

// Highly nested structure for version control and recording
const recordManager = {
    versionMap: {},
    recordUpdate(id, val, optName, upValue, isRemote, rid) {
        const lamportTime = cClock.localEvent();
        this.incrementVersion(id);
        const logEntry = formatLog(id, val, optName, upValue, lamportTime, isRemote, rid);
        this.writeLogToFile(id, logEntry);
        if (!isRemote) {
            updateList.push(new Update(id, id2ver[id], val));
        }
    },
    incrementVersion(id) {
        if (!this.versionMap[id]) {
            this.versionMap[id] = 0;
        }
        this.versionMap[id]++;
    },
    writeLogToFile(id, logEntry) {
        const filePath = path.join(__dirname, 'DBs', `${id}.txt`);
        fs.appendFileSync(filePath, logEntry);
    },
    getVersion(id) {
        return this.versionMap[id] || 0;
    }
};

// Separate class to handle JSON marshalling/unmarshalling
class JSONHandler {
    static marshal(data) {
        return JSON.stringify(data);
    }

    static unmarshal(jsonStr) {
        return JSON.parse(jsonStr);
    }
}

// Complex Set Manager which depends on several other classes
class SetManager {
    constructor() {
        this.sets = {};
    }

    createSet(id) {
        const newSet = new Set(id);
        this.sets[id] = newSet;
        return newSet;
    }

    getSet(id) {
        return this.sets[id];
    }

    handleSetOperation(id, operation, value) {
        const set = this.getSet(id);
        if (!set) throw new Error(`Set with id ${id} not found!`);
        set.SetVal(operation, value);
    }
}

class Set {
    constructor(id) {
        this.id = id;
        this.values = {};
        this.updates = [];
        this.addhoc_reject_local_update = {};
        createFile(id);
        this.recordManager = recordManager;
    }

    Id() {
        return this.id;
    }

    Values() {
        return this.values;
    }

    Contain(value) {
        return this.values.hasOwnProperty(value);
    }

    Add(value) {
        if (!this.Contain(value)) {
            this.values[value] = true;
        }
    }

    Remove(value) {
        if (this.Contain(value)) {
            delete this.values[value];
        }
    }

    SetVal(opt_name, value) {
        if (opt_name === 'Add') {
            this.Add(value);
        } else if (opt_name === 'Remove') {
            this.Remove(value);
        }
        const cur_upte = [opt_name, value.toString()];
        this.updates.push(cur_upte);

        if (opt_name === 'Add') {
            let rejct_up = `${opt_name}_${value}`;
            this.addhoc_reject_local_update[rejct_up] = true;
        }
        const val = addhoc_keys4Persist(this.Values());
        this.recordManager.recordUpdate(this.Id(), val, opt_name, value, false);
    }

    SetRemoteVal(rid, opt_name, value) {
        if (opt_name === 'Add') {
            this.Add(value);
        } else if (opt_name === 'Remove') {
            this.Remove(value);
        }
        this.recordManager.recordUpdate(this.Id(), this.Values(), opt_name, value, true, rid);
    }

    Union(o) {
        const m = {};
        Object.keys(this.values).forEach((item) => m[item] = true);
        Object.keys(o.values).forEach((item) => {
            if (!m.hasOwnProperty(item)) {
                this.values[item] = true;
            }
        });
    }

    Merge(rid, r_updates) {
        const tmpSet = new Set(rid);
        r_updates.forEach(([opt_name, val]) => {
            let req_val = parseInt(val, 10);
            if (opt_name === 'Add') {
                tmpSet.Add(req_val);
            } else {
                tmpSet.Remove(req_val);
            }
        });

        this.Union(tmpSet);
        this.Print();
    }

    Print() {
        const values = Object.keys(this.values).map(Number).sort((a, b) => a - b);
        console.log(`Set:${this.id}`, values);
        this.recordManager.recordUpdate(this.id, JSON.stringify(values), "Print", 0, false);
    }

    ToMarshal() {
        this.updates.push([this.Id().toString(), ""]);
        return JSONHandler.marshal(this.updates);
    }

    addhoc_setRemoteWrap(rid, opt_name, value) {
        this.recordManager.recordUpdate(this.Id(), JSONHandler.marshal(this.Values()), opt_name, value, true, rid);
    }

    addhoc_mergeWrap(rid, r_updates) {
        if (r_updates.length > 0 && r_updates[0][0] === 'undo') {
            let [undo_opt, undo_val, remote_Lt] = r_updates[0][1].split('_');
            if (undo(rid, undo_opt, parseInt(undo_val, 10), parseInt(remote_Lt, 10))) {
                if (undo_opt === 'Add') {
                    this.SetRemoteVal(rid, 'Remove', parseInt(undo_val, 10));
                    this.addhoc_setRemoteWrap(rid, 'Remove', parseInt(undo_val, 10));
                } else if (undo_opt === 'Remove') {
                    this.SetRemoteVal(rid, 'Add', parseInt(undo_val, 10));
                    this.addhoc_setRemoteWrap(rid, 'Add', parseInt(undo_val, 10));
                }
            }
        } else {
            this.Merge(rid, r_updates);
        }
    }
}

// Creates a file for storing logs related to a specific Set
function createFile(id) {
    const dirPath = path.join(__dirname, 'DBs');
    if (!fs.existsSync(dirPath)) {
        fs.mkdirSync(dirPath);
    }
    const filePath = path.join(dirPath, `${id}.txt`);
    fs.writeFileSync(filePath, "");
}

// Special Update class for record management
class Update {
    constructor(id, version, value) {
        this.ID = id;
        this.Version = version;
        this.Value = value;
    }
}

// Highly nested function that generates and persists keys
function addhoc_keys4Persist(values) {
    const sortedKeys = Object.keys(values).map(Number).sort((a, b) => a - b);
    const jsonStr = JSONHandler.marshal(sortedKeys);
    return jsonStr;
}

// Lamport Clock logic with complex dependency on Set and record keeping
class LamportClock {
    constructor() {
        this.time = 0;
        this.mu = false;  // Simulating a mutex lock
    }

    localEvent() {
        this.incrementTime();
        return this.time;
    }

    incrementTime() {
        if (!this.mu) {
            this.mu = true; // Lock
            this.time++;
            this.mu = false; // Unlock
        }
    }

    getTime() {
        return this.time;
    }
}

// Clock instances, update tracking, and record keeping
const cClock = new LamportClock();
const rid2LT = {};
const id2ver = {};
const updateList = [];

// Undo function that depends on multiple modules
function undo(rID, undoUpdate, undoVal, rLT) {
    if (rid2LT[rID].getTime() > rLT) {
        return true;
    }
    return false;
}
