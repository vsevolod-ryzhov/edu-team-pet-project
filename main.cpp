#include <iostream>
#include <queue>
#include <vector>
#include <algorithm>

using namespace std;

int main() {
    queue<pair<int, int>> q;
    queue<pair<int, int>> res;
    int W, H;
    string temp;
    vector<string> matrix;
    cin >> W >> H;
    // input and checking
    while(cin >> temp) {
        if (temp.size() < W || temp.size() > W) {
            cout << "Error size" << endl;
            return 0;
        }
        matrix.push_back(temp);
    }
    if (matrix.size() != H) {
        cout << "Error size" << endl;
        return 0;
    }
    pair<int, int> check = {-1, -1};
    pair<int, int> START = {-1, -1};
    pair<int, int> END = {-1, -1};
    for (int i = 0; i < H; i++) {
        if (matrix[i].find('S') != -1) {
            if (START != check) {
                START = {i + 1, matrix[i].find('S')};
            }
            else {
                cout << "The most of Start points" << endl;
                return 0;
            }
        }
        if (matrix[i].find('G') != -1) {
            if (END != check) {
                END = {i + 1, matrix[i].find('S')};
            }
            else {
                cout << "The most of Finish points" << endl;
                return 0;
            }
        }
    }
    if (END == check || START == check) {
        cout << "No Finish point or Start point" << endl;
        return 0;
    }


    return 0;
}