using System;
using System.Collections.Generic;
using System.IO;
using System.Runtime.InteropServices;
using System.Linq;
using RAC.Payloads;
using static RAC.Errors.Log;

namespace RAC.Operations
{
    public static class Name
    {
        public static string ReplicaName;
    }

    public class Interceptor : PNCounter
    {
        public Interceptor(string uid, Parameters parameters) : base(uid, parameters)
        {
        }

        public override Responses SetValue()
        {
            PNCPayload pl = new PNCPayload(uid, (int)Config.numReplicas, (int)Config.replicaId);
            Name.ReplicaName = pl.uid;
            Persist_Helper.Create_File(Name.ReplicaName);
            int value = this.parameters.GetParam<int>(0);
            Persist_Helper.Record(Name.ReplicaName, value.ToString());
            return base.SetValue();
        }


        public Responses Reverse()
        {

            Responses res = new Responses(Status.success);
            PNCPayload pl = new PNCPayload(uid, (int)Config.numReplicas, (int)Config.replicaId);

            IntPtr ptr = Persist_Helper.undo(Name.ReplicaName, -2);
            if (ptr != IntPtr.Zero)
            {
                // string rev = Marshal.PtrToStringAnsi(ptr);
                string rev = Marshal.PtrToStringUTF8(ptr);
                Console.WriteLine("rev is " + rev);
                if (rev != null)
                {
                    try
                    {
                        int value = int.Parse(rev);
                        Console.WriteLine("undo requires");
                        Console.WriteLine("rollback value is " + value);

                        if (value >= 0)
                            pl.PVector[pl.replicaid] = value;
                        else
                            pl.NVector[pl.replicaid] = -value;

                        //pl.PVector[pl.replicaid] = 0;
                        //pl.NVector[pl.replicaid] = 0;

                        Persist_Helper.Record(Name.ReplicaName, value.ToString());

                        this.payload = pl;
                    }
                    catch (FormatException)
                    {
                        Console.WriteLine("Invalid format. Please enter a valid string.");
                    }
                }

            }
            else
            {
                // Console.WriteLine("An exception occurred");
            }


            base.GenerateSyncRes(ref res);
            res.AddResponse(Dest.client);
            return res;
        }

        public Responses Increment()
        {
            this.payload.PVector[this.payload.replicaid] += this.parameters.GetParam<int>(0);

            // ADDED

            int value = this.payload.PVector.Sum() - this.payload.NVector.Sum();
            int delta = this.parameters.GetParam<int>(0);
            Persist_Helper.Record(Name.ReplicaName, value.ToString());

            this.payload.PVector[this.payload.replicaid] -= this.parameters.GetParam<int>(0);
            return base.Increment();

        }

        public Responses Decrement()
        {
            this.payload.NVector[this.payload.replicaid] += this.parameters.GetParam<int>(0);
            // ADDED
            int value = this.payload.PVector.Sum() - this.payload.NVector.Sum();
            int delta = this.parameters.GetParam<int>(0) * (-1);
            Persist_Helper.Record(Name.ReplicaName, value.ToString());


            this.payload.NVector[this.payload.replicaid] -= this.parameters.GetParam<int>(0);
            return base.Decrement();
        }

        public override Responses Synchronization()
        {
            return base.Synchronization();
        }
    }

}

