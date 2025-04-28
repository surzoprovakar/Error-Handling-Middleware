using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using RAC.Payloads;
using static RAC.Errors.Log;

namespace RAC.Operations
{
    /// <summary>
    /// Reversible Graph
    /// </summary>
    public class Graph : Operation<RGraphPayload>
    {
        public static class Name
        {
            public static string ReplicaName;
        }
        private const int TAG_LEN = 8;

        // todo: set this to its typecode
        public override string typecode { get; set; } = "gh";

        public Graph(string uid, Parameters parameters) : base(uid, parameters)
        {
            // todo: put any necessary data here
        }

        public string Process(string inputs)
        {
            string res;

            string[] lines = inputs.Split('\n');
            // string vtx = lines[2].Replace(" ", "");
            // string edgs = lines[4].Replace("|", "").Replace(" ", "");

            // Console.WriteLine("vtx " + vtx);
            // Console.WriteLine("edgs " + edgs);

            // res = vtx + "" + edgs;
            // Console.WriteLine("res " + res);
            string vtx = lines[2];
            string edgs = lines[4];

            res = vtx + " <||||> " + edgs;
            // Maybe try after encrypting and decrypting this value
            // Problem is the white space and the special characters
            return res;
        }
        public override Responses GetValue()
        {

            Responses res = new Responses(Status.success);

            StringBuilder sb = new StringBuilder();

            sb.Append("Vertices:\n");

            foreach (var v in this.payload.vertices)
                sb.Append(v.value + " ");

            sb.Append("\nEdges:\n");

            foreach (var e in this.payload.edges)
            {
                string v1 = e.Item1.v1;
                string v2 = e.Item1.v2;

                if (lookup(v1) == (null, null) || lookup(v2) == (null, null))
                    continue;
                else
                    sb.Append("(" + v1 + "," + v2 + ")|");
            }

            res.AddResponse(Dest.client, sb.ToString());
            return res;
        }

        public override Responses SetValue()
        {
            this.payload = new RGraphPayload(uid);
            Name.ReplicaName = this.payload.uid;
            Persist_Helper.Create_File(Name.ReplicaName);

            this.payload = new RGraphPayload(uid);

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

            Responses res = new Responses(Status.success);
            res.AddResponse(Dest.client);

            GenerateSyncRes(ref res, "av", v.ToString());

            if (string.Join("", res.contents).Contains("Succeed"))
            {
                string inputs = string.Join("", base.GetValue().contents);

                string vals = Process(inputs);
                Persist_Helper.Record(Name.ReplicaName, vals);

            }

            return res;
        }

        public Responses RemoveVertex()
        {
            string value = this.parameters.GetParam<string>(0);

            Responses res;

            (string, string) toRemove;

            // precondition
            // collect all unique pairs in V containing v
            if ((toRemove = lookup(value)) == (null, null))
            {
                res = new Responses(Status.fail);
                res.AddResponse(Dest.client, "Vertex DNE");
                return res;
            }

            // v is not the head of an existing arc
            foreach (var item in this.payload.edges)
            {
                if (item.Item1.v1 == value)
                {
                    res = new Responses(Status.fail);
                    res.AddResponse(Dest.client, "Vertex is the head of an existing arc");
                    return res;
                }
            }

            // effect (R)
            this.payload.vertices.Remove(toRemove);
            res = new Responses(Status.success);
            res.AddResponse(Dest.client);

            GenerateSyncRes(ref res, "rv", toRemove.ToString());

            if (string.Join("", res.contents).Contains("Succeed"))
            {
                string inputs = string.Join("", base.GetValue().contents);

                string vals = Process(inputs);
                Persist_Helper.Record(Name.ReplicaName, vals);
            }

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

            // A := A ∪ {((v′, v′′),w)}
            string tag = UniqueTag();
            var e = ((v1, v2), tag);
            this.payload.edges.Add(e);

            res = new Responses(Status.success);
            GenerateSyncRes(ref res, "ae", e.ToString());
            res.AddResponse(Dest.client);

            if (string.Join("", res.contents).Contains("Succeed"))
            {
                string inputs = string.Join("", base.GetValue().contents);

                string vals = Process(inputs);
                Persist_Helper.Record(Name.ReplicaName, vals);

            }

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

            // A := A \ R
            this.payload.edges.Remove(toRemove);
            res = new Responses(Status.success);
            GenerateSyncRes(ref res, "re", toRemove.ToString());
            res.AddResponse(Dest.client);

            if (string.Join("", res.contents).Contains("Succeed"))
            {
                string inputs = string.Join("", base.GetValue().contents);
                string vals = Process(inputs);
                Persist_Helper.Record(Name.ReplicaName, vals);

            }

            return res;

        }

        // look up vertex
        private (string, string) lookup(string vertex)
        {

            foreach (var item in this.payload.vertices)
            {
                if (item.value == vertex)
                    return item;
            }

            return (null, null);
        }

        // look up edge
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
                    break;
                case "av":
                    var v = (updateSplit[0], updateSplit[1]);
                    this.payload.vertices.Add(v);
                    break;
                case "rv":
                    var vremove = (updateSplit[0], updateSplit[1]);
                    this.payload.vertices.Remove(vremove);
                    break;
                case "ae":
                    var e = ((updateSplit[0], updateSplit[1]), updateSplit[2]);
                    this.payload.edges.Add(e);
                    break;
                case "re":
                    var eremove = ((updateSplit[0], updateSplit[1]), updateSplit[2]);
                    this.payload.edges.Remove(eremove);
                    break;
            }


            return new Responses(Status.success);
        }

        public void GenerateSyncRes(ref Responses res, string type, string update)
        {
            Parameters syncPm = new Parameters(2);

            // type: 
            // "n": new graph
            // "av": add vertex
            // "rv": remove vertex
            // "ae": add edge
            // "re": remove edge
            syncPm.AddParam(0, type);
            // effect-update msg
            syncPm.AddParam(1, update);

            string broadcast = Parser.BuildCommand(this.typecode, "y", this.uid, syncPm);
            res.AddResponse(Dest.broadcast, broadcast, false);
        }

        // TODO: move this to utli class
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

        static string[] ExtractValues(string input, string pattern)
        {
            var match = Regex.Match(input, pattern);
            if (match.Success)
            {
                string valuesString = match.Groups[1].Value;
                return valuesString.Split(' ');
            }
            return new string[0]; // Return an empty array if no match is found
        }

        // ToDo : Need to optimize more
        static List<Tuple<string, string>> ExtractTuples(string input, string pattern)
        {
            var match = Regex.Match(input, pattern);
            List<Tuple<string, string>> tuples = new List<Tuple<string, string>>();

            if (match.Success)
            {
                string tuplesString = match.Groups[1].Value;
                string[] tupleStrings = tuplesString.Split('|');

                foreach (var tupleString in tupleStrings)
                {
                    string[] elements = tupleString.Trim('(', ')').Split(',');
                    if (elements.Length == 2)
                    {
                        tuples.Add(new Tuple<string, string>(elements[0].Trim(), elements[1].Trim()));
                    }
                }
            }

            return tuples;
        }


        // ToDo : Need to optimize more
        public Responses Reverse()
        {
            IntPtr ptr = Persist_Helper.undo(Name.ReplicaName, 1);
            if (ptr != IntPtr.Zero)
            {
                string rev = Marshal.PtrToStringAnsi(ptr);
                Console.WriteLine("rev is " + rev);
                if (rev != null)
                {
                    try
                    {
                        Console.WriteLine("undo requires");
                        //int value = int.Parse(rev);
                        //Console.WriteLine("rollback value is " + rev);
                        int line = int.Parse(rev);
                        string fileName = "src/Operations/DBs/" + Name.ReplicaName + ".txt";
                        string roll_back_val = File.ReadLines(fileName).Skip(line).FirstOrDefault();

                        Console.WriteLine("rollback value is " + roll_back_val);

                        string[] vector1 = ExtractValues(roll_back_val, @Name.ReplicaName + " " + (line + 1).ToString() + " (.*?) <||||>");
                        Console.WriteLine("vertices: " + string.Join(" ", vector1));
                        List<Tuple<string, string>> vector2 = ExtractTuples(roll_back_val, @"<\|\|\|\|> (.*?)$");

                        base.payload.vertices.Clear();
                        base.payload.edges.Clear();

                        foreach (var v in vector1)
                        {
                            string tag = base.UniqueTag();
                            var v1 = (v, tag);
                            base.payload.vertices.Add(v1);
                        }

                        // HashSet<string> vtxHashSet = new HashSet<string>(base.payload.vertices.Select(tuple => tuple.Item1.ToString()));
                        // string vtx = string.Join("|", vtxHashSet);
                        // Console.WriteLine("testing: " + vtx);

                        foreach (var ed in vector2)
                        {
                            //Console.WriteLine("itm edge is: " + ed);
                            string tag = base.UniqueTag();
                            var e = ((ed.Item1, ed.Item2), tag);
                            base.payload.edges.Add(e);
                        }

                        string inputs = string.Join("", base.GetValue().contents);

                        string vals = Process(inputs);
                        Persist_Helper.Record(Name.ReplicaName, vals);
                    }
                    catch (FormatException)
                    {
                        Console.WriteLine("Invalid format. Please enter a valid string.");
                    }
                }
                else
                {
                    // Console.WriteLine("An exception occurred");
                }
            }

            var res = new Responses(Status.success);
            res.AddResponse(Dest.client);
            return res;
        }

    }

    class ConstCharPtrMarshaler : ICustomMarshaler
    {
        public object MarshalNativeToManaged(IntPtr pNativeData)
        {
            return Marshal.PtrToStringAnsi(pNativeData);
        }

        public IntPtr MarshalManagedToNative(object ManagedObj)
        {
            return IntPtr.Zero;
        }

        public void CleanUpNativeData(IntPtr pNativeData)
        {
        }

        public void CleanUpManagedData(object ManagedObj)
        {
        }

        public int GetNativeDataSize()
        {
            return IntPtr.Size;
        }

        static readonly ConstCharPtrMarshaler instance = new ConstCharPtrMarshaler();

        public static ICustomMarshaler GetInstance(string cookie)
        {
            return instance;
        }
    }
    public static class Persist_Helper
    {
        private static Dictionary<string, int> id2ver = new Dictionary<string, int>();
        private static List<Tuple<string, int, string>> updateList = new List<Tuple<string, int, string>>();
        private static string dbPath = "src/Operations/DBs/";

        // Create the file for a given replica
        public static void CreateFile(string id)
        {
            if (!Directory.Exists(dbPath))
            {
                Directory.CreateDirectory(dbPath);
            }

            string filename = Path.Combine(dbPath, $"{id}.txt");
            using (StreamWriter sw = File.AppendText(filename))
            {
                Console.WriteLine($"Storage file created for Replica {id}");
            }

            if (!id2ver.ContainsKey(id))
            {
                id2ver[id] = 0;
            }
        }

        // Record the new value in the file and update the in-memory structure
        public static void Record(string id, string val)
        {
            string filename = Path.Combine(dbPath, $"{id}.txt");
            id2ver[id]++;

            using (StreamWriter sw = new StreamWriter(filename, append: true))
            {
                sw.WriteLine($"{id} {id2ver[id]} {val}");
            }

            Tuple<string, int, string> entry = new Tuple<string, int, string>(id, id2ver[id], val);
            updateList.Add(entry);
        }

        // Process return value to extract the rollback element
        public static string ProcessRet(string retVal)
        {
            string rollElement = retVal.Substring(retVal.LastIndexOf(" ") + 1);
            return rollElement.Trim();
        }

        // Process values to return last N updates in an int array
        public static int[] ProcessVals(int prevUpdates)
        {
            int[] lastValues = new int[prevUpdates];
            int index = 0;
            foreach (var entry in updateList.AsEnumerable().Reverse())
            {
                string strVal = entry.Item3;
                if (index < prevUpdates)
                {
                    lastValues[index] = int.Parse(strVal);
                    index++;
                }
                if (index >= prevUpdates)
                {
                    break;
                }
            }
            return lastValues;
        }

        // Undo function based on provided options
        public static string Undo(string id, int optNums)
        {
            string filename = Path.Combine(dbPath, $"{id}.txt");
            string retVal;

            // Undo based on data corruption detection
            if (optNums == -2)
            {
                if (DbDataCorruptUndo(3)) // Hardcoded example with 3 updates, can be changed
                {
                    retVal = MultipleUndo(filename, 3);
                    return ProcessRet(retVal);
                }
                else
                {
                    Console.WriteLine("No undo required based on data corruption.");
                    return null;
                }
            }

            // Undo single update operation
            if (optNums == 1)
            {
                Console.WriteLine("Undoing the last update operation.");
                retVal = updateList[updateList.Count - 2].Item3;
            }
            else
            {
                Console.WriteLine($"Undoing the last {optNums} update operations.");
                retVal = updateList[updateList.Count - (optNums + 1)].Item3;
            }

            return ProcessRet(retVal);
        }

        // Mock function for multiple undo (replace with real implementation)
        private static string MultipleUndo(string filename, int optNums)
        {
            // Mock implementation for multiple undo
            return updateList[updateList.Count - optNums].Item3;
        }

        // Check if data is corrupted and should perform undo
        public static bool DbDataCorruptUndo(int prevUpdates)
        {
            // In C#, we don’t have dlopen/dlsym like in C++, instead we use DllImport for interop
            // Assuming the SHA256 plugin functionality is loaded through another managed library
            // For this example, we'll mock the function
            int[] data = ProcessVals(prevUpdates);

            return ShouldUndo(data, prevUpdates);
        }

        // Mock function for checking corruption, replace with actual plugin logic
        private static bool ShouldUndo(int[] data, int prevUpdates)
        {
            // Simulating a corruption check (e.g., checksum mismatch)
            // Replace with real logic later
            return data.Sum() % 2 == 0; // Example check: undo if sum of data is even
        }
    }

}

