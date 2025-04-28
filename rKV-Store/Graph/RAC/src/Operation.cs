#define EAGER
#undef EAGER

using RAC.Payloads;
using RAC.Errors;
using static RAC.Errors.Log;
using RAC.History;

namespace RAC.Operations
{   
    /// <summary>
    /// Base class for all CRDT replication algorithms. 
    /// Also holds infomation on current operation.
    /// </summary>
    /// <typeparam name="PayloadType">Payload class of this CRDT</typeparam>
    public abstract class Operation<PayloadType> where PayloadType: Payload
    {
        public string uid { get; }
        public abstract string typecode { get; set; }
        public Parameters parameters { protected set; get; }
        
#if EAGER
        public OpHistoryEager history;
#else 
        public OpHistory history;
#endif

        /// <summary>
        /// Payload that actually holds the states of a CRDT.
        /// Other methods should use this attribute to gain access to the
        /// states of CRDT and set this attribtue when state is changed/updated.
        /// </summary>
        public PayloadType payload { protected set; get; }

        /// <summary>
        /// Used by Save(), if this set to true, the this.payload will not
        /// be saved by the Memory Manager for this operation.
        /// </summary>
        protected bool noSideEffect = false;

        /// <summary>
        /// Constructor:
        /// </summary>
        /// <param name="uid">uid of the accessing object for this op</param>
        /// <param name="parameters">any parameters passed in for this op</param>
        public Operation(string uid, Parameters parameters)
        {
            this.uid = uid;
            this.parameters = parameters;

            try
            {
                this.payload = (PayloadType) Global.memoryManager.GetPayload(uid);
            }
            catch (PayloadNotFoundException) 
            {
                this.payload = null;
            }

#if EAGER
            try
            {
                this.history = (OpHistoryEager)Global.memoryManager.history[uid];
                if (!(this is HistoryHandler))
                    history.Compensate = this.Compensate;
            }
            catch (System.Collections.Generic.KeyNotFoundException)
            {
                Global.memoryManager.history.Add(uid, new OpHistoryEager(uid, this.Compensate));
                this.history = (OpHistoryEager)Global.memoryManager.history[uid];

            }
#else 
            try
            {
                this.history = Global.memoryManager.history[uid];
            }
            catch (System.Collections.Generic.KeyNotFoundException)
            {
                Global.memoryManager.history.Add(uid, new OpHistory(uid, this.Compensate));
                this.history = Global.memoryManager.history[uid];

            }
#endif

            
        }

        /// <summary>
        /// Call the memory manager to store the this.payload. 
        /// </summary>
        public void Save()
        {
            if (!noSideEffect)
            {
                Global.memoryManager.StorePayload(uid, payload);
                DEBUG(uid + " successfully stored");
            }
        }

        // =================================!IMPORTANT!===========================================
        // ====!!NEXT 3 METHODS MUST BE IMPLEMENTED AND PROVIDED AS APIs FOR ANY CRDT!!======

        // Request Handling APIs:
        /// <summary>
        /// Set the given value in this.parameters to CRDT object with the given uid.
        /// If uid DNE, create a new CRDT object.
        /// This method must be provided by any CRDT.
        /// </summary>
        /// <returns>A Responses instance containing any information to the client or
        /// broadcasting messages.</returns>
        public abstract Responses SetValue();

        /// <summary>
        /// Get the value of a CRDT object with given uid.
        /// This method must be provided by any CRDT.
        /// </summary>
        /// <returns>A Responses instance containing any information to the client or
        /// broadcasting messages.</returns>
        public abstract Responses GetValue();

        /// <summary>
        /// Synchronize the state of the object based on the the given values in 
        /// this.paramters. Normally this is a request come from other replicas.
        /// The content request could either be states or delta-state for state-based
        /// CRDT, or effect-update for operation-based CRDT.
        /// This method must be provided by any CRDT.
        /// </summary>
        /// <returns>A Responses instance containing any information to the client or
        /// broadcasting messages.</returns>
        public abstract Responses Synchronization();

        // ========================================================================================

        /// <summary>
        ///  
        /// </summary>
        /// <returns>A Responses instance containing any information to the client or
        /// broadcasting messages.</returns>
        public virtual Responses DeleteValue()
        {
            Responses res = new Responses(Status.success);

            // TODO: deletion things

            return res;
        }

        /// <summary>
        /// Used for reverse
        /// </summary>
        public virtual void Compensate(string opid)
        {
            
        }

    }
    
}