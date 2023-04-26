package dbus;

import org.freedesktop.dbus.interfaces.DBusInterface;
import org.freedesktop.dbus.types.UInt32;
import com.sun.jna.Library;
import com.sun.jna.Native;
import org.freedesktop.dbus.connections.impl.DBusConnection;
import org.freedesktop.dbus.exceptions.DBusException;
import org.freedesktop.dbus.interfaces.ObjectManager;
import org.freedesktop.dbus.messages.*;

import java.io.FileNotFoundException;
import java.io.InvalidObjectException;
import java.net.ConnectException;
import java.nio.Buffer;
import java.nio.ByteBuffer;
import java.nio.CharBuffer;

//import org.freedesktop.dbus.test.helper.signals.SampleSignals.TestSignal;
import org.freedesktop.dbus.InternalSignal;

public class Smb2 {
    private String addressWithPort;
    private String user;
    private String password;
    private String shareName;
    private GoSmb smb;
    private final String errStringDetection = "#ERROR#";
    private String libPath = "../../bin/lib-smb2.so";

    private interface GoSmb extends Library {

        // ------------------------------------------------------------------------------------
        // CLIENT SECTION
        // ------------------------------------------------------------------------------------

        String InitClient(String addressWithPort, String user, String psw, String shareName)
                throws UnsupportedOperationException;

        boolean IsConnected();

        void CloseConn();

        String ListShares();

        // ------------------------------------------------------------------------------------
        // DIRECTORY SECTION
        // ------------------------------------------------------------------------------------
        String CreateFolder(String name) throws UnsupportedOperationException;

        String DeleteFolder(String name) throws UnsupportedOperationException;

        String RenameFolder(String oldPath, String newPath) throws UnsupportedOperationException;

        String CheckIfFolderExists(String name);

        String IsDir(String name);

        // ------------------------------------------------------------------------------------
        // FILE SECTION
        // ------------------------------------------------------------------------------------
        String AppendLine(String fileName, String strToWrite);

        String WriteStringFromOffset(String fileName, String strToWrite, long offset);

        String ReadFile(String string);

        String ReadFileWithOffsets(String filename, long offset, long chunkSize);

        String RemoveFile(String fileName);

        String RenameFile(String oldPath, String newPath);

    }

    /**
     * Constructor
     */
    public Smb2(String addressWithPort, String user, String password, String shareName) {
        this.addressWithPort = addressWithPort;
        this.user = user;
        this.password = password;
        this.shareName = shareName;
        String lib = libPath;
        smb = (GoSmb) Native.loadLibrary(lib, GoSmb.class);
    }

    /**
     * Constructor with possibility to connect automatically
     */
    public Smb2(String addressWithPort, String user, String password, String shareName, boolean connect)
            throws ConnectException {
        new Smb2(addressWithPort, user, password, shareName);
        String lib = libPath;
        smb = (GoSmb) Native.loadLibrary(lib, GoSmb.class);
        if (connect) {
            connect();
        }
    }

    /**
     * Init samba2 client to establish the connection
     * 
     * @return true if OK! Otherwise exception
     */
    public boolean connect() throws ConnectException {
        String err = smb.InitClient(this.addressWithPort, this.user, this.password, this.shareName);
        if (err == null || err.isEmpty()) {
            return true;
        }
        throw new ConnectException(err);
    }

    /**
     * Close connection
     */
    public void disconnect() {
        smb.CloseConn();

    }

    /**
     * Check if is connected
     * 
     * @return boolean
     */
    public boolean isConnected() {
        return smb.IsConnected();
    }

    /**
     * List all shares
     * 
     * @return String with names separated by comma es: myshare1,myshare2. If string
     *         contains #ERROR# an error occured
     */
    public String[] listShares() throws FileNotFoundException {
        String errOrList = smb.ListShares();
        System.out.println(errOrList);
        if (errOrList.contains(errStringDetection)) {
            throw new FileNotFoundException(errOrList);
        }
        return errOrList.split(",");
    }

    /**
     * Append line to the end of file
     * 
     * @param filename is considered as path
     * @param text     is the the string to append
     * @return boolean
     * @exception if file not found
     */
    public boolean appendLine(String filename, String text) throws FileNotFoundException {
        String err = smb.AppendLine(filename, text);
        if (err != null) {
            throw new FileNotFoundException(err);
        }
        return true;
    }

    /**
     * Write a string starting from a specific offset
     * 
     * @param filename is considered as path
     * @param text     is the the string to write
     * @param offset   where to start
     * @return boolean
     * @exception if file not found
     */
    public boolean writeStringFromOffset(String filename, String text, long offset) throws FileNotFoundException {
        String err = smb.WriteStringFromOffset(filename, text, offset);
        if (err != null) {
            throw new FileNotFoundException(err);
        }
        return true;
    }

    /**
     * Read All File
     * 
     * @param filename is considered as path
     * @return file as string
     * @exception if file not found
     */
    public String readFile(String filename) throws FileNotFoundException {
        String errOrFile = smb.ReadFile(filename);
        if (errOrFile != null && errOrFile.contains(errStringDetection)) {
            throw new FileNotFoundException(errOrFile);
        }
        return errOrFile;
    }

    /**
     * Delete File
     * 
     * @param filename is considered as path
     * @return file as string
     * @exception if file not found
     */
    public boolean removeFile(String filename) throws FileNotFoundException {
        String errOrFile = smb.RemoveFile(filename);
        if (errOrFile != null && errOrFile.contains(errStringDetection)) {
            throw new FileNotFoundException(errOrFile);
        }
        return true;
    }

    /**
     * Create Folder
     * 
     * @param name is consideredas path
     * @exception if file name or path is incorrect
     */
    public String createFolder(String name) throws InvalidObjectException {
        String err = smb.CreateFolder(name);
        if (err != null && err.contains(errStringDetection)) {
            throw new InvalidObjectException(err);
        }
        return err;
    }

    /**
     * Rename Folder (move)
     * 
     * @param oldPath is considered as path
     * @param newPath is considered as path
     * @return if null or empty is OK! otherwise error occurs
     * @exception if path is not valid
     */
    public String renameFolder(String oldPath, String newPath) throws InvalidObjectException {
        String err = smb.RenameFolder(oldPath, newPath);
        if (err != null && err.contains(errStringDetection)) {
            throw new InvalidObjectException(err);
        }
        return err;
    }

    /**
     * Delete Folder
     * 
     * @param name is considered as path
     * @return if null or empty is OK! otherwise error occurs
     * @exception if path is not valid
     */
    public String deleteFolder(String name) throws InvalidObjectException {
        String err = smb.DeleteFolder(name);
        if (err != null && err.contains(errStringDetection)) {
            throw new InvalidObjectException(err);
        }
        return err;
    }

    /**
     * Check If Folder Exists
     * 
     * @param name is considered as path
     * @return boolean
     * @exception if path is not valid
     */
    public boolean checkIfFolderExists(String name) throws InvalidObjectException {
        String errOrExists = smb.CheckIfFolderExists(name);
        if (errOrExists != null && errOrExists.contains(errStringDetection)) {
            if (errOrExists.contains("file does not exist")) {
                return false;
            }
            throw new InvalidObjectException(errOrExists);
        }
        System.out.println(name + " exists?" + errOrExists);
        boolean retB = Boolean.parseBoolean(errOrExists);
        return retB;
    }

    /**
     * Check if a path is a dir, otherwise means file
     * 
     * @param name is considered as path
     * @return boolean
     * @exception if path is not valid
     */
    public boolean isDir(String name) throws FileNotFoundException {
        String errOrExists = smb.IsDir(name);
        if (errOrExists != null && errOrExists.contains(errStringDetection)) {
            throw new FileNotFoundException(errOrExists);
        }
        boolean retB = Boolean.parseBoolean(errOrExists);
        return retB;
    }

}