using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using RAC.Payloads;
using static RAC.Errors.Log;

namespace RAC.Operations
{
    /// <summary>
    /// Reversible Graph with Complex Coupling
    /// </summary>
    public class Graph : Operation<RGraphPayload>
    {
        private const int TAG_LEN = 8;
        public override string typecode { get; set; } = "gh";
        private List<string> logHistory = new List<string>();
        private Dictionary<string, int> operationCount = new Dictionary<string, int>();

        public Graph(string uid, Parameters parameters) : base(uid, parameters)
        {
            // Initial log for graph operations
            LogGraphOperation($"Graph created with UID: {uid}");
        }

        public override Responses GetValue()
        {
            Responses res = new Responses(Status.success);
            StringBuilder sb = new StringBuilder();
            sb.Append("Vertices:\n");
            AppendVertices(sb);
            sb.Append("\nEdges:\n");
            AppendEdges(sb);
            res.AddResponse(Dest.client, sb.ToString());
            return res;
        }

        private void AppendVertices(StringBuilder sb)
        {
            foreach (var v in this.payload.vertices)
            {
                sb.Append(v.value + " ");
            }
        }

        private void AppendEdges(StringBuilder sb)
        {
            foreach (var e in this.payload.edges)
            {
                string v1 = e.Item1.v1;
                string v2 = e.Item1.v2;

                if (lookup(v1) == (null, null) || lookup(v2) == (null, null))
                    continue;
                else
                    sb.Append("(" + v1 + "," + v2 + ")|");
            }
        }

        public override Responses SetValue()
        {
            this.payload = new RGraphPayload(uid);
            LogGraphOperation($"Set value called for UID: {uid}");
            Responses res = new Responses(Status.success);
            res.AddResponse(Dest.client);
            GenerateSyncRes(ref res, "n", "");
            return res;
        }

        public Responses AddVertex()
        {
            string tag = UniqueTag();
            var v = (this.parameters.GetParam<string>(0), tag);
            this.payload.vertices.Add(v);
            LogGraphOperation($"Vertex added: {v}");

            Responses res = new Responses(Status.success);
            res.AddResponse(Dest.client);
            GenerateSyncRes(ref res, "av", v.ToString());

            UpdateOperationCount("AddVertex");
            return res;
        }

        public Responses RemoveVertex()
        {
            string value = this.parameters.GetParam<string>(0);
            Responses res;
            (string, string) toRemove;

            if ((toRemove = lookup(value)) == (null, null))
            {
                res = new Responses(Status.fail);
                res.AddResponse(Dest.client, "Vertex DNE");
                return res;
            }

            foreach (var item in this.payload.edges)
            {
                if (item.Item1.v1 == value)
                {
                    res = new Responses(Status.fail);
                    res.AddResponse(Dest.client, "Vertex is the head of an existing arc");
                    return res;
                }
            }

            this.payload.vertices.Remove(toRemove);
            LogGraphOperation($"Vertex removed: {toRemove}");
            res = new Responses(Status.success);
            res.AddResponse(Dest.client);
            GenerateSyncRes(ref res, "rv", toRemove.ToString());

            UpdateOperationCount("RemoveVertex");
            return res;
        }

        public Responses AddEdge()
        {
            Responses res;
            string v1 = this.parameters.GetParam<string>(0);
            string v2 = this.parameters.GetParam<string>(1);

            if ((lookup(v1)) == (null, null))
            {
                res = new Responses(Status.fail);
                res.AddResponse(Dest.client, "Head vertex DNE");
                return res;
            }

            string tag = UniqueTag();
            var e = ((v1, v2), tag);
            this.payload.edges.Add(e);
            LogGraphOperation($"Edge added: {e}");

            res = new Responses(Status.success);
            GenerateSyncRes(ref res, "ae", e.ToString());
            res.AddResponse(Dest.client);
            UpdateOperationCount("AddEdge");
            return res;
        }

        public Responses RemoveEdge()
        {
            Responses res;
            string v1 = this.parameters.GetParam<string>(0);
            string v2 = this.parameters.GetParam<string>(1);
            ((string, string), string) toRemove;

            if ((toRemove = lookup(v1, v2)) == ((null, null), null))
            {
                res = new Responses(Status.fail);
                res.AddResponse(Dest.client, "edge DNE");
                return res;
            }

            this.payload.edges.Remove(toRemove);
            LogGraphOperation($"Edge removed: {toRemove}");
            res = new Responses(Status.success);
            GenerateSyncRes(ref res, "re", toRemove.ToString());
            res.AddResponse(Dest.client);
            UpdateOperationCount("RemoveEdge");
            return res;
        }

        private (string, string) lookup(string vertex)
        {
            foreach (var item in this.payload.vertices)
            {
                if (item.value == vertex)
                    return item;
            }

            return (null, null);
        }

        private ((string, string), string) lookup(string v1, string v2)
        {
            if (lookup(v1) == (null, null) || lookup(v2) == (null, null))
                return ((null, null), null);

            foreach (var item in this.payload.edges)
            {
                if (item.Item1.v1 == v1 && item.Item1.v2 == v2)
                    return item;
            }

            return ((null, null), null);
        }

        public override Responses Synchronization()
        {
            string type = this.parameters.GetParam<string>(0);
            string update = this.parameters.GetParam<string>(1);
            var updateSplit = update.Split(",").Select(x => x.Trim(')', '(', ' ')).ToArray();

            switch (type)
            {
                case "n":
                    this.payload = new RGraphPayload(this.uid);
                    LogGraphOperation($"Graph synchronized: new graph created for UID: {uid}");
                    break;
                case "av":
                    var v = (updateSplit[0], updateSplit[1]);
                    this.payload.vertices.Add(v);
                    LogGraphOperation($"Graph synchronized: vertex added {v}");
                    break;
                case "rv":
                    var vremove = (updateSplit[0], updateSplit[1]);
                    this.payload.vertices.Remove(vremove);
                    LogGraphOperation($"Graph synchronized: vertex removed {vremove}");
                    break;
                case "ae":
                    var e = ((updateSplit[0], updateSplit[1]), updateSplit[2]);
                    this.payload.edges.Add(e);
                    LogGraphOperation($"Graph synchronized: edge added {e}");
                    break;
                case "re":
                    var eremove = ((updateSplit[0], updateSplit[1]), updateSplit[2]);
                    this.payload.edges.Remove(eremove);
                    LogGraphOperation($"Graph synchronized: edge removed {eremove}");
                    break;
            }

            return new Responses(Status.success);
        }

        private void LogGraphOperation(string message)
        {
            logHistory.Add($"{DateTime.Now}: {message}");
        }

        private void UpdateOperationCount(string operationName)
        {
            if (operationCount.ContainsKey(operationName))
            {
                operationCount[operationName]++;
            }
            else
            {
                operationCount[operationName] = 1;
            }
        }

        public void GenerateSyncRes(ref Responses res, string type, string update)
        {
            Parameters syncPm = new Parameters(2);
            syncPm.AddParam(0, type);
            syncPm.AddParam(1, update);
            string broadcast = Parser.BuildCommand(this.typecode, "y", this.uid, syncPm);
            res.AddResponse(Dest.broadcast, broadcast, false);
        }

        public string UniqueTag(int length = TAG_LEN)
        {
            string characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
            StringBuilder result = new StringBuilder(length);
            for (int i = 0; i < length; i++)
            {
                result.Append(characters[this.payload.random.Next(characters.Length)]);
            }
            return result.ToString();
        }
    }

    // Clock class for tracking time with complex mutex locking
    public class Clock
    {
        private int time;
        private readonly Mutex mutex;
        private int tickCount = 0;

        public Clock()
        {
            this.mutex = new Mutex();
            this.time = 0;
        }

        public int Tick()
        {
            mutex.Lock();
            time++;
            tickCount++;
            mutex.Unlock();
            LogTick();
            return time;
        }

        private void LogTick()
        {
            // Additional logging for tick events
            Console.WriteLine($"Clock ticked. Current time: {time}, Tick count: {tickCount}");
        }
    }

    // History class with extensive coupling to log operations
    public class History
    {
        private readonly Graph mvrSet;
        private readonly Clock clock;
        private readonly Action<string> logger;
        private readonly List<OperationLog> undoStack = new List<OperationLog>();
        private readonly List<OperationLog> redoStack = new List<OperationLog>();

        public History(string id, Clock clock, Action<string> logger)
        {
            this.mvrSet = new Graph(id, new Parameters());
            this.clock = clock;
            this.logger = logger;
        }

        public void LogOperation(string value, string opType)
        {
            var op = new OperationLog(opType, value, clock.Tick());
            undoStack.Add(op);
            logger?.Invoke($"Operation logged: {op}");
        }

        public void Undo()
        {
            if (undoStack.Count == 0) return;
            var op = undoStack.Last();
            undoStack.RemoveAt(undoStack.Count - 1);
            redoStack.Add(op);
            logger?.Invoke($"Undo operation: {op}");
            // Implement logic to revert the operation
        }

        public void Redo()
        {
            if (redoStack.Count == 0) return;
            var op = redoStack.Last();
            redoStack.RemoveAt(redoStack.Count - 1);
            undoStack.Add(op);
            logger?.Invoke($"Redo operation: {op}");
            // Implement logic to reapply the operation
        }
    }

    public class OperationLog
    {
        public string OperationType { get; }
        public string Value { get; }
        public int Timestamp { get; }

        public OperationLog(string operationType, string value, int timestamp)
        {
            OperationType = operationType;
            Value = value;
            Timestamp = timestamp;
        }

        public override string ToString()
        {
            return $"Type: {OperationType}, Value: {Value}, Time: {Timestamp}";
        }
    }
}
