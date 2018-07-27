import json
import argparse
#jsonlog_path = "c:\\Users\\user\\Desktop\\Workspace\\python\\.vscode\\samplelog.txt"
# python logdup.py -h 
# gives help for using script
if __name__ == '__main__':
    #Parse argument for script
    parser = argparse.ArgumentParser(description='Check json log file for duplicate userId(/appId)')
    parser.add_argument('jsonlog_path',help="path to log file")
    parser.add_argument('-app', action='store_true',help='use appId for key?')
    
    args = parser.parse_args()
    jsonlog_path = args.jsonlog_path
    app = args.app
    #print(jsonlog_path)
    log_file = open(jsonlog_path)

    # parse log line by line
    userId_appId_list = {}
    duplicate_userId_list = {}
    log_list = log_file.readlines()
    for single_log in log_list : 
        log_dump = json.loads(single_log)
        if app :
            key_in_log_dump = set(('userId' ,'appId')) <= set(log_dump)
        else:
            key_in_log_dump = 'userId' in log_dump

        if key_in_log_dump and 'region' in log_dump['body'] :
            cur_userId = log_dump['userId']
            if app:
                cur_appId = log_dump['appId']
                key = (cur_userId,cur_appId)
            else:
                key = cur_userId
            
            cur_region = log_dump['body']['region']
            if key in userId_appId_list and userId_appId_list[key] != cur_region :
                if key not in duplicate_userId_list :
                    duplicate_userId_list.setdefault(key,[userId_appId_list[key]])
                duplicate_userId_list[key].append(cur_region)
            else :
                userId_appId_list[key] = cur_region

    # output result
    '''
    print("Format 1")
    print(duplicate_userId_list)
    print("\nFormat 2")
    '''
    for dupkey in duplicate_userId_list:
        print(dupkey)
    
