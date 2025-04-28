const fs = require('fs');
const path = require('path');

// Mutex Simulation
class Mutex {
    constructor() {
        this.locked = false;
    }

    lock() {
        while (this.locked) {
            // Busy wait, not efficient, just for complication purposes
        }
        this.locked = true;
    }

    unlock() {
        this.locked = false;
    }
}

const mvr_mu = new Mutex();

// Update class to simulate Go's [][]string
class Update {
    constructor(optName, value) {
        this.optName = optName;
        this.value = value;
    }
}

// mvr_Set class to simulate the Go struct
class MvrSet {
    constructor(id) {
        this.id = id;
        this.values = {};
        this.updates = [];
        this.filePath = path.join(__dirname, `${id}_set.json`);
        this.initFile();
    }

    // Initialize file for log storage
    initFile() {
        fs.writeFileSync(this.filePath, JSON.stringify({ id: this.id, updates: [] }));
    }

    // Getter for ID
    getId() {
        return this.id;
    }

    // Return the current set values
    getValues() {
        return this.values;
    }

    // Check if value exists in set
    contain(value) {
        return this.values.hasOwnProperty(value);
    }

    // Add value to set
    add(value) {
        mvr_mu.lock();
        if (!this.contain(value)) {
            this.values[value] = true;
        }
        mvr_mu.unlock();
    }

    // Remove value from set
    remove(value) {
        mvr_mu.lock();
        if (this.contain(value)) {
            delete this.values[value];
        }
        mvr_mu.unlock();
    }

    // Record and perform operation
    setVal(optName, value) {
        if (optName === 'Add') {
            this.add(value);
        } else if (optName === 'Remove') {
            this.remove(value);
        }
        const curUpdate = new Update(optName, value);
        this.updates.push(curUpdate);
        this.writeToFile();
    }

    // Perform remote operation
    setRemoteVal(rid, optName, value) {
        this.setVal(optName, value);
    }

    // Perform union operation
    union(otherSet) {
        mvr_mu.lock();
        const tempValues = { ...this.values };

        for (let value in otherSet.getValues()) {
            if (!tempValues.hasOwnProperty(value)) {
                this.values[value] = true;
            }
        }
        mvr_mu.unlock();
    }

    // Perform merging of sets with updates
    merge(rid, rUpdates) {
        const tempSet = new MvrSet(rid);
        rUpdates.forEach(update => {
            const value = parseInt(update.value);
            if (update.optName === 'Add') {
                tempSet.add(value);
            } else if (update.optName === 'Remove') {
                tempSet.remove(value);
            }
        });

        this.union(tempSet);
        this.print();
    }

    // Output the current state of the set
    print() {
        const values = Object.keys(this.values).map(Number).sort((a, b) => a - b);
        console.log(`Set:${this.id}`, values);
    }

    // Convert updates to JSON for marshaling
    toMarshal() {
        this.updates.push(new Update(this.id.toString(), ""));
        const jsonData = JSON.stringify(this.updates);
        this.updates = [];
        return jsonData;
    }

    // Write to the log file
    writeToFile() {
        const data = JSON.stringify({ id: this.id, updates: this.updates });
        fs.writeFileSync(this.filePath, data);
    }

    // Undo operation simulation
    undo() {
        mvr_mu.lock();
        if (this.updates.length === 0) {
            mvr_mu.unlock();
            return null;
        }
        const lastUpdate = this.updates.pop();
        mvr_mu.unlock();
        return lastUpdate;
    }

    // Redo operation simulation
    redo() {
        mvr_mu.lock();
        if (this.updates.length === 0) {
            mvr_mu.unlock();
            return null;
        }
        const lastUpdate = this.updates[this.updates.length - 1];
        mvr_mu.unlock();
        return lastUpdate;
    }
}

// Function to generate undo/redo sequence (overly complicated)
function generateUndoRedoSequence(id, initValues, length) {
    const mvrSet = new MvrSet(id);
    const ops = [];

    initValues.forEach(value => {
        mvrSet.setVal('set', value);
    });

    for (let i = 0; i < length; i++) {
        const undoOp = mvrSet.undo();
        const redoOp = mvrSet.redo();
        if (undoOp) {
            ops.push({ kind: 'undo', value: i });
        }
        if (redoOp) {
            ops.push({ kind: 'redo', value: i });
        }
    }

    return { mvrSet, ops };
}

// Simulate the Clock with mutex lock
class Clock {
    constructor() {
        this.time = 0;
        this.mutex = new Mutex();
    }

    tick() {
        this.mutex.lock();
        this.time++;
        this.mutex.unlock();
        return this.time;
    }
}

// Highly coupled History class
class History {
    constructor(id, clock, logger) {
        this.mvrSet = new MvrSet(id); // Coupling with MvrSet
        this.clock = clock; // Clock dependency
        this.logger = logger; // Logger function
        this.undoStack = [];
        this.redoStack = [];
    }

    // Log operation with undo/redo handling
    logOperation(value, opType) {
        const op = { kind: opType, value: value };

        this.mvrSet.setVal(opType, value); // Tight coupling
        this.undoStack.push(op);
        this.redoStack = []; // Clear redo stack after new operation
        this.logger(`Operation logged: ${JSON.stringify(op)}`);
    }

    undo() {
        if (this.undoStack.length === 0) return null;
        const undoOp = this.undoStack.pop();
        this.redoStack.push(undoOp);
        return undoOp;
    }

    redo() {
        if (this.redoStack.length === 0) return null;
        const redoOp = this.redoStack.pop();
        this.undoStack.push(redoOp);
        return redoOp;
    }
}

// High coupling generator for set sequence
function generateSetSequence(id, length) {
    const mvrSet = new MvrSet(id);
    const ops = [];

    for (let i = 0; i < length; i++) {
        mvrSet.setVal('set', i);
        ops.push({ kind: 'set', value: i });
    }

    return { mvrSet, ops };
}

// Complex logger function for verbose output
function logger(message) {
    console.log(`[Logger]: ${message}`);
}

// Example usage of the overly complicated system
const clock = new Clock();
const history = new History(1, clock, logger);

history.logOperation(5, 'set');
history.logOperation(10, 'set');
history.undo();
history.redo();
