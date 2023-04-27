package dbus;

import static org.junit.Assert.assertEquals;
import static org.junit.jupiter.api.Assertions.*;

import org.freedesktop.dbus.exceptions.DBusException;
import org.junit.jupiter.api.Test;

class TestGo {

    Smb2 smb = new Smb2("xxx:xxx:xxx:xxx:pppp", "user", "psw", "shareName");

    @Test
    void test() {
        try {

            smb.connect();
            String[] list = smb.listShares();
            smb.appendLine("provona.txt", "ciao");
            smb.writeStringFromOffset("provona.txt", "$333$$", 2000);
            smb.appendLine("toDel.txt", "");
            smb.removeFile("toDel.txt");
            System.out.println(list.length);
            // smb.createFolder("babbaluba");
            smb.disconnect();
        } catch (Exception e) {
            // TODO Auto-generated catch block
            System.out.println("errore da go " + e.getMessage());
            e.printStackTrace();
        } finally {

        }
    }

    @Test
    void writeFiletest() {
        try {
            smb.connect();
            double randomNumber = Math.random();
            smb.appendLine(randomNumber + "test.txt", "" + randomNumber);
            String res = smb.readFile(randomNumber + "test.txt");
            boolean bRes = res.equalsIgnoreCase("" + randomNumber);
            assertEquals("check", bRes, true);
        } catch (Exception ex) {

        } finally {
            smb.disconnect();
        }
    }

    @Test
    void foldersTest() {
        try {
            smb.connect();
            double randomNumber = Math.random();
            String folderName = (randomNumber + "test").replace(".", "C");
            smb.createFolder(folderName);
            boolean ok = smb.checkIfFolderExists(folderName);

            assertEquals("check", ok, true);
            smb.renameFolder(folderName, folderName + "_new");
            assertEquals("check", smb.checkIfFolderExists(folderName + "_new"), true);
            ok = smb.checkIfFolderExists(folderName);
            assertEquals("check", ok, false);

            smb.deleteFolder(folderName + "_new");
            assertEquals("check", smb.checkIfFolderExists(folderName + "_new"), false);
        } catch (Exception ex) {
            ex.printStackTrace();
        } finally {
            smb.disconnect();
        }
    }

    @Test
    void isConnectedTest() {
        try {
            smb.connect();
            assertEquals("check", smb.isConnected(), true);
            smb.disconnect();
            assertEquals("check", smb.isConnected(), false);
        } catch (Exception ex) {
            ex.printStackTrace();
        } finally {

        }
    }

}
