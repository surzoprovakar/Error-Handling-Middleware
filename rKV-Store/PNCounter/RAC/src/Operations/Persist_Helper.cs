// Example.cs

using System;
using System.Runtime.InteropServices;
using System.Text;

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
    public const string so_path = "/home/provakar/Oxygen/Distributed-Undo-Redo/rKV-Store/Graph/RAC/src/Operations/libPersist.so";
    [DllImport(so_path, CharSet = CharSet.Ansi, CallingConvention = CallingConvention.Cdecl)]
    public static extern void Create_File([MarshalAs(UnmanagedType.LPStr)] string id);

    [DllImport(so_path, CharSet = CharSet.Ansi, CallingConvention = CallingConvention.Cdecl)]
    public static extern void Record([MarshalAs(UnmanagedType.LPStr)] string id, [MarshalAs(UnmanagedType.LPStr)] string val);

    // [DllImport("/home/provakar/Oxygen/Undo-As-Service/C#/rKVCRDT-main/RAC/src/Operations/libPersist.so", CharSet = CharSet.Ansi, CallingConvention = CallingConvention.Cdecl)]
    // //[return: MarshalAs(UnmanagedType.CustomMarshaler, MarshalTypeRef = typeof(ConstCharPtrMarshaler))]
    // [return: MarshalAs(UnmanagedType.BStr)]
    // public static extern StringBuilder undo([MarshalAs(UnmanagedType.LPStr)] string id);

    [DllImport(so_path, CharSet = CharSet.Ansi, CallingConvention = CallingConvention.Cdecl)]
    public static extern IntPtr undo([MarshalAs(UnmanagedType.LPStr)] string id, int opt_nums);
}