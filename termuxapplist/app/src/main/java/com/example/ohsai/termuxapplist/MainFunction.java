package com.example.ohsai.termuxapplist;

import android.content.Context;
import android.content.Intent;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.util.Log;

import java.io.File;
import java.io.FileOutputStream;
import java.util.List;

public class MainFunction {
    private static final String TAG = MainActivity.class.getSimpleName();
    public static void appListRecreate(Context context) {
        final PackageManager pm = context.getPackageManager();
//get a list of installed apps.
        List<ApplicationInfo> packages = pm.getInstalledApplications(PackageManager.GET_META_DATA);

        String datapath = "/storage/emulated/0/TermuxApps/";
        File directory = new File(datapath);
        if (! directory.exists()){
            directory.mkdir();
        }
        if (directory.isDirectory())
        {
            String[] children = directory.list();
            for (int i = 0; i < children.length; i++)
            {
                new File(directory, children[i]).delete();
            }
        }
        for (ApplicationInfo packageInfo : packages) {
            Intent temp = pm.getLaunchIntentForPackage(packageInfo.packageName);
            if(temp != null){ //for every launchable application
                //Log for debugging purpose
                Log.d(TAG, "Installed package :" + packageInfo.packageName);
                Log.d(TAG, "Source dir : " + packageInfo.sourceDir);
                String appfullname = temp.getComponent().toShortString();
                Log.d(TAG, "PackageName/ActivityName :" + appfullname);
                if(packageInfo.packageName.equals("com.google.android.youtube")) {
                    appfullname = "{com.google.android.youtube/com.google.android.youtube.HomeActivity}";
                }
                //create file
                String filename = pm.getApplicationLabel(packageInfo).toString();
                //String filename = "Testing";
                FileOutputStream outputStream;
                File file = new File(datapath + filename);

                //write output
                String output = "#!/system/bin/sh \nam start -n "+appfullname
                        .replace("{","").replace("}","");
                try {
                    outputStream = new FileOutputStream(file);
                    outputStream.write(output.getBytes());
                    outputStream.close();
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
        }
    }
}
