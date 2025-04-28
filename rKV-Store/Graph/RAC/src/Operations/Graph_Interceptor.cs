using System;
using System.IO;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Runtime.InteropServices;
using RAC.Payloads;
using static RAC.Errors.Log;
using System.Text.RegularExpressions;

namespace RAC.Operations
{
    /// <summary>
    /// Reversible Graph
    /// </summary>
    public static class Name
    {
        public static string ReplicaName;
    }
    public class Graph_Interceptor : Graph
    {
        private const int TAG_LEN = 8;

        // todo: set this to its typecode

        public Graph_Interceptor(string uid, Parameters parameters) : base(uid, parameters)
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

        public override Responses SetValue()
        {
            this.payload = new RGraphPayload(uid);
            Name.ReplicaName = this.payload.uid;
            Persist_Helper.Create_File(Name.ReplicaName);
            return base.SetValue();

        }

        public Responses AddVertex()
        {
            Responses res = base.AddVertex();
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

            Responses res = base.RemoveVertex();

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

            Responses res = base.AddEdge();

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

            Responses res = base.RemoveEdge();
            if (string.Join("", res.contents).Contains("Succeed"))
            {
                string inputs = string.Join("", base.GetValue().contents);
                string vals = Process(inputs);
                Persist_Helper.Record(Name.ReplicaName, vals);

            }

            return res;

        }

        // ToDo : Need to optimize more
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

        // public string UniqueTag(int length = TAG_LEN)
        // {
        //     string characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
        //     StringBuilder result = new StringBuilder(length);
        //     for (int i = 0; i < length; i++)
        //     {
        //         result.Append(characters[this.payload.random.Next(characters.Length)]);
        //     }
        //     return result.ToString();
        // }

    }

}

