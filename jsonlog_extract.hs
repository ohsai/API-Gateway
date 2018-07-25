import System.Environment (getArgs)
import Data.List.Split
import Data.List
main = do
    args <- getArgs
    content <- readFile (args !! 0)
    let linesOfFiles = lines content
    let temp = process_input linesOfFiles 
    mapM_ print $ create_duplist temp
    

duplicate x y =  (key x == key y) && not (region x == region y) 
    where 
        key = filter (\(x:xs) -> (x=="appId" || x == "userId" ))
        region = filter (\(x:xs) -> (x=="region" ))

exists_duplicate x xlist = find (duplicate x) xlist 

create_duplist xlist = [x | Just x <- maybe_duplist]
    where
        maybe_duplist = map (\obj -> exists_duplicate obj xlist) xlist
        
jsonfield input = map (\xs -> (splitOn ":" [ x | x <- xs, not (x `elem` "{}\"") ])) input

process_input input = map (filter (\(x:xs) -> (x=="appId" || x == "userId" || x == "region"))) ( map jsonfield (map (splitOn ",") input)) 
