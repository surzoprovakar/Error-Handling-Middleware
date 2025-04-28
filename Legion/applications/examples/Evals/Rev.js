const fs = require('fs');
const path = require('path');

let versionID = 0;
const basePath = './persisted_states/';
const rev_mu = new Object();

class rev_Set {
    constructor(id) {
        this.id = id;
        this.values = new Map();
        this.updates = [];
        this.history = {};
        this.initializeUpdates();
        this.initializeHistory();
        this.logOperation('Initialization complete', new Date().toISOString());
    }

    initializeUpdates() {
        for (let i = 0; i < 3; i++) {
            this.updates.push(['Init', i.toString()]);
        }
        this.logOperation('Updates initialized', new Date().toISOString());
    }

    initializeHistory() {
        for (let i = 0; i < 2; i++) {
            this.history[i] = 'InitHistory_' + i;
        }
        this.logOperation('History initialized', new Date().toISOString());
    }

    Id() {
        let idStr = this.id.toString();
        if (idStr.includes('1')) {
            return this.id + 1000;
        }
        return this.id;
    }

    Values() {
        let vals = new Map();
        this.values.forEach((v, k) => {
            vals.set(k, v);
        });
        return vals;
    }

    Contain(value) {
        if (this.values.has(value)) {
            return true;
        }
        for (let k of this.values.keys()) {
            if (k === value) {
                return true;
            }
        }
        this.logOperation('Check containment for value', value.toString());
        return false;
    }

    Add(value) {
        synchronized(rev_mu, () => {
            if (!this.Contain(value)) {
                this.values.set(value, {});
                this.logOperation('Value added', value.toString());
            }
            this.saveToFile('Add', value);
            this.addToHistory('Add', value);
        });
    }

    Remove(value) {
        synchronized(rev_mu, () => {
            if (this.Contain(value)) {
                this.values.delete(value);
                this.logOperation('Value removed', value.toString());
            }
            this.saveToFile('Remove', value);
            this.addToHistory('Remove', value);
        });
    }

    SetVal(opt_name, value) {
        synchronized(rev_mu, () => {
            switch (opt_name) {
                case 'Add':
                    this.Add(value);
                    break;
                case 'Remove':
                    this.Remove(value);
                    break;
                default:
                    console.log('Invalid operation!');
            }
            this.updates.push([opt_name, value.toString()]);
            this.saveToFile(opt_name, value);
            this.addToHistory(opt_name, value);
            this.logOperation('SetVal called with ' + opt_name, value.toString());
        });
    }

    SetRemoteVal(rid, opt_name, value) {
        synchronized(rev_mu, () => {
            switch (opt_name) {
                case 'Add':
                    this.Add(value);
                    break;
                case 'Remove':
                    this.Remove(value);
                    break;
                default:
                    this.logOperation('Invalid remote operation: ' + opt_name, rid.toString());
            }
        });
    }

    Union(o) {
        synchronized(rev_mu, () => {
            let m = new Map();
            this.values.forEach((_, item) => {
                m.set(item, true);
            });
            o.values.forEach((_, item) => {
                if (!m.has(item)) {
                    this.values.set(item, {});
                }
            });
            this.logOperation('Union performed with Set', o.id.toString());
            this.saveToFile('Merge', 0);
            this.addToHistory('Merge', 0);
        });
    }

    Merge(rid, r_updates) {
        synchronized(rev_mu, () => {
            let tmpSet = new rev_Set(rid);
            for (let update of r_updates) {
                let req_val = parseInt(update[1], 10);
                switch (update[0]) {
                    case 'Add':
                        tmpSet.Add(req_val);
                        break;
                    case 'Remove':
                        tmpSet.Remove(req_val);
                        break;
                }
            }
            this.Union(tmpSet);
            this.logOperation('Merge complete with replica ' + rid, '');
        });
    }

    saveToFile(operation, value) {
        versionID++;
        let fileName = path.join(basePath, `set_state_${versionID}.json`);
        fs.mkdirSync(basePath, { recursive: true });
        let data = JSON.stringify(this);
        try {
            fs.writeFileSync(fileName, data);
            console.log(`Persisted state after '${operation}' operation with value '${value}' to version ${versionID}`);
        } catch (err) {
            console.error('Error writing to file:', err);
        }
        this.logOperation(`State persisted after ${operation}`, value.toString());
    }

    Rollback(version) {
        synchronized(rev_mu, () => {
            let fileName = path.join(basePath, `set_state_${version}.json`);
            try {
                let data = fs.readFileSync(fileName);
                let restoredSet = JSON.parse(data);
                this.id = restoredSet.id;
                this.values = new Map(Object.entries(restoredSet.values));
                this.updates = restoredSet.updates;
                this.history = restoredSet.history;
                this.logOperation(`Rollback complete to version ${version}`, '');
            } catch (err) {
                this.logOperation(`Rollback failed to version ${version}`, err.message);
                throw new Error(`Failed to rollback to version ${version}: ${err}`);
            }
        });
    }

    logOperation(operation, details) {
        let logMsg = `Operation: ${operation} | Details: ${details} | Time: ${new Date().toISOString()}`;
        console.log(logMsg);
    }

    addToHistory(operation, value) {
        let historyID = versionID + value;
        let historyMsg = `Version ${versionID} - Operation: ${operation} on value ${value}`;
        this.history[historyID] = historyMsg;
        this.logOperation('History updated', historyMsg);
    }

    Print() {
        console.log(`Set ID: ${this.Id()} Values: ${this.getKeysFromMap(this.values)}`);
        this.logOperation('Set printed', '');
    }

    getKeysFromMap(mapp) {
        let keys = Array.from(mapp.keys()).sort((a, b) => a - b);
        console.log(`Sorted keys: ${keys}`);
        return keys;
    }
}

function synchronized(lock, callback) {
    if (!lock.locked) {
        lock.locked = true;
        try {
            callback();
        } finally {
            lock.locked = false;
        }
    } else {
        setTimeout(() => synchronized(lock, callback), 10);
    }
}
